package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/entity"
	applog "github.com/ktigay/goph-keeper/internal/log"
	"github.com/ktigay/goph-keeper/internal/server/config"
	appdb "github.com/ktigay/goph-keeper/internal/server/db"
	datahandler "github.com/ktigay/goph-keeper/internal/server/handler/grpc"
	"github.com/ktigay/goph-keeper/internal/server/interceptor"
	userrepo "github.com/ktigay/goph-keeper/internal/server/repository/user"
	userdatarepo "github.com/ktigay/goph-keeper/internal/server/repository/userdata"
	"github.com/ktigay/goph-keeper/internal/server/security"
	authsrv "github.com/ktigay/goph-keeper/internal/server/service/auth"
	userdatasrv "github.com/ktigay/goph-keeper/internal/server/service/userdata"
)

func main() {
	ctx := context.TODO()

	var (
		cfg    *config.Config
		logger *slog.Logger
		pool   *pgxpool.Pool
		err    error
	)

	if cfg, err = config.New(os.Args[1:]); err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	logger = applog.New(cfg.LogLevel, os.Stdout)
	logger.Debug("config loaded", "config", cfg)

	if pool, err = appdb.NewPgxPool(ctx, cfg.DatabaseDSN); err != nil {
		log.Fatalf("Failed to create connect to DB: %v", err)
	}
	if err = appdb.CreateSchema(ctx, pool); err != nil {
		log.Fatalf("Failed to create structure: %v", err)
	}

	var (
		jwtAuth = security.NewJWTWrapper[entity.Identity](cfg.AuthSecret)

		// txFacade  = appdb.NewPgxTxFacade(pool)
		dbWrapper = appdb.NewTxConnWrapper(pool)

		userRepo = userrepo.New(dbWrapper, logger)
		userSrv  = authsrv.New(userRepo)

		userdataRepo = userdatarepo.New(dbWrapper, logger)
		userdataSrv  = userdatasrv.New(userdataRepo)
	)

	exitCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.WithRecover(logger),
			interceptor.WithLogging(logger),
			interceptor.NewAuth(jwtAuth, interceptor.AccessList()).WithAuthorization(),
		),
	)
	data.RegisterUserDataServiceServer(grpcServer, datahandler.NewUserDataHandler(userdataSrv))
	auth.RegisterAuthServiceServer(grpcServer, datahandler.NewAuthHandler(userSrv, jwtAuth))
	reflection.Register(grpcServer)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var listen net.Listener
		listen, err = net.Listen("tcp", cfg.ServerGRPCHost)
		if err != nil {
			log.Fatalf("can't listen: %v", err)
		}

		logger.Debug("grpc server listening", "address", listen.Addr().String())

		if err = grpcServer.Serve(listen); err != nil {
			log.Fatalf("can't start grpc server: %v", err)
		}
		wg.Done()
	}()

	go func() {
		<-exitCtx.Done()

		if grpcServer != nil {
			grpcServer.GracefulStop()
			logger.Debug("grpc server gracefully stopped")
		}
	}()

	wg.Wait()

	logger.Debug("server shutdown gracefully")
}
