package app

import (
	"context"
	"log/slog"

	"github.com/rivo/tview"

	"github.com/ktigay/goph-keeper/internal/client/entity"
	authhandler "github.com/ktigay/goph-keeper/internal/client/tui/handler/auth"
	userdatahanler "github.com/ktigay/goph-keeper/internal/client/tui/handler/userdata"
	apppage "github.com/ktigay/goph-keeper/internal/client/tui/page"
	"github.com/ktigay/goph-keeper/internal/client/tui/page/auth"
	"github.com/ktigay/goph-keeper/internal/client/tui/page/userdatalist"
	e "github.com/ktigay/goph-keeper/internal/entity"
)

// Api api сервисы.
type Api struct {
	AuthSrv         authhandler.Service
	UserDataSrv     userdatahanler.Service
	UserDataSyncSrv authhandler.SyncService
}

// New создаёт консольное приложение.
func New(ctx context.Context, api Api, logger *slog.Logger, isSyncedCh <-chan bool, signedInCh, quitCh chan<- struct{}) *tview.Application {
	app := tview.NewApplication()
	appPages := apppage.NewPages()

	loginHandler := authhandler.New(api.AuthSrv, api.UserDataSyncSrv)
	loginView := auth.New(
		auth.Callbacks{
			OnSignIn: func(credentials entity.Credentials) error {
				err := loginHandler.SignIn(ctx, credentials)
				if err != nil {
					return err
				}
				appPages.SwitchToPage(apppage.UserDataList)

				close(signedInCh)

				logger.Debug("user signed in")
				return nil
			},
			OnSignUp: func(credentials entity.Credentials) error {
				err := loginHandler.SignUp(ctx, credentials)
				if err != nil {
					logger.Debug("user signed up with", "error", err.Error())
					return err
				}

				logger.Debug("user signed up")
				return nil
			},
		},
	)

	userDataHandler := userdatahanler.New(api.UserDataSrv)
	var userDataView *userdatalist.Page
	userDataView = userdatalist.New(
		userdatalist.Callbacks{
			OnItemUpdate: func(data e.UserData) error {
				_, err := userDataHandler.ItemUpdate(ctx, data)
				return err
			},
			OnItemAdd: func(data e.UserData) error {
				_, err := userDataHandler.ItemAdd(ctx, data)
				return err
			},
			OnRefreshData: func() {
				err := api.UserDataSyncSrv.SyncFromRemote(ctx)
				if err != nil {
					logger.Debug("user data sync complete with error", "error", err.Error())
				}
				userDataView.Render()
			},
			OnItemDelete: func(data e.UserData) error {
				return userDataHandler.ItemDelete(ctx, data.UUID)
			},
			OnQuit: func() {
				close(quitCh)
			},
		},
		func() ([]e.UserData, error) {
			return userDataHandler.GetList(ctx)
		},
	)

	appPages.AddPage(apppage.Auth, loginView, true, true)
	appPages.AddPage(apppage.UserDataList, userDataView, true, false)

	go func() {
		for {
			select {
			case <-isSyncedCh:
				userDataView.Render()
			case <-ctx.Done():
				return
			}
		}
	}()

	app.SetRoot(appPages.TviewPages, true).EnableMouse(true)
	return app
}
