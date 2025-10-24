package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authclient "github.com/ktigay/goph-keeper/internal/client/client/grpc/auth"
	userdataclient "github.com/ktigay/goph-keeper/internal/client/client/grpc/userdata"
	"github.com/ktigay/goph-keeper/internal/client/config"
	"github.com/ktigay/goph-keeper/internal/client/interceptor"
	authrepo "github.com/ktigay/goph-keeper/internal/client/repository/auth"
	userdatarepo "github.com/ktigay/goph-keeper/internal/client/repository/userdata"
	authsrv "github.com/ktigay/goph-keeper/internal/client/service/auth"
	syncsrv "github.com/ktigay/goph-keeper/internal/client/service/sync"
	userdatasrv "github.com/ktigay/goph-keeper/internal/client/service/userdata"
	"github.com/ktigay/goph-keeper/internal/client/tui/app"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/entity"
	applog "github.com/ktigay/goph-keeper/internal/log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	var (
		loc        *time.Location
		cfg        *config.Config
		fileLogger *os.File
		logger     *slog.Logger
		grpcClient *grpc.ClientConn
		err        error
	)
	loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc

	ctx := context.Background()

	if cfg, err = config.New(os.Args[1:]); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	if cfg.Version {
		buildInfo()
		os.Exit(0)
	}

	if fileLogger, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666); err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer func() {
		if err = fileLogger.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}()
	logger = applog.New(cfg.LogLevel, fileLogger)

	var (
		authRepo       *authrepo.Repository
		userDataRepo   *userdatarepo.Repository
		authClient     *authclient.Client
		userDataClient *userdataclient.Client
		authSrv        *authsrv.Service
		userDataSrv    *userdatasrv.Service
		syncSrv        *syncsrv.Service
	)

	authRepo = authrepo.New()
	grpcClient, err = grpc.NewClient(
		cfg.ServerGRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.TimeoutInterceptor(cfg.SrvRequestTimeout),
			interceptor.AuthInterceptor(authRepo),
		),
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}

	authClient = authclient.New(auth.NewAuthServiceClient(grpcClient))
	userDataClient = userdataclient.New(data.NewUserDataServiceClient(grpcClient))
	authSrv = authsrv.New(authClient, authRepo, logger)

	userDataRepo = userdatarepo.New()
	userDataSrv = userdatasrv.New(userDataRepo)
	syncSrv = syncsrv.New(userDataClient, userDataRepo, logger)

	signedInCh := make(chan struct{})

	isSyncedCh := make(chan bool, 1)
	defer close(isSyncedCh)

	quitCh := make(chan struct{})

	exitCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	consoleApp := app.New(exitCtx, app.Api{
		AuthSrv:         authSrv,
		UserDataSrv:     userDataSrv,
		UserDataSyncSrv: syncSrv,
	}, logger, isSyncedCh, signedInCh, quitCh)

	wg := &sync.WaitGroup{}

	go func() {
		<-signedInCh
		wg.Add(1)

		ticker := time.NewTicker(time.Duration(cfg.SrvSyncToInterval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				var d []entity.UserData
				if d, err = syncSrv.SyncToRemote(ctx); err != nil {
					logger.Debug("SyncToRemote error", "error", err)
				}
				if len(d) > 0 {
					isSyncedCh <- true
				}
			case <-exitCtx.Done():
				logger.Debug("SyncToRemote exit")
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		if err = consoleApp.Run(); err != nil {
			panic(err)
		}
		logger.Debug("consoleApp exited")
		wg.Done()
		stop()
	}()

	go func() {
		<-exitCtx.Done()
		consoleApp.Stop()
	}()

	go func() {
		<-quitCh
		stop()
	}()

	wg.Wait()

	logger.Debug("client shutdown gracefully")
}

func buildInfo() {
	_, _ = fmt.Fprintf(os.Stdout, `  Build version: %s
  Build date: %s
`, buildVersion, buildDate)
}
