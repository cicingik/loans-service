// Package app ...
package app

import (
	"context"

	"github.com/cicingik/loans-service/config"
	"github.com/cicingik/loans-service/delivery"
	"github.com/cicingik/loans-service/repository/auth"
	funding "github.com/cicingik/loans-service/repository/loan_funding"
	"github.com/cicingik/loans-service/repository/loans"
	"github.com/cicingik/loans-service/repository/postgre"
	"github.com/cicingik/loans-service/repository/users"
	"github.com/cicingik/loans-service/services"
	"go.uber.org/fx"
)

// WebApplication ...
type WebApplication struct {
	*fx.App
	cfg *config.AppConfig
	db  *postgre.DbEngine
}

// NewWebApplication ...
func NewWebApplication(cfg *config.AppConfig) (*WebApplication, error) {
	app := &WebApplication{}

	container := fx.New(
		fx.Provide(
			func() *config.AppConfig {
				return cfg
			},
			NewDatabase,
			NewHTTPServer,
		),

		fx.Provide(auth.NewAuthRepository),
		fx.Provide(users.NewUsersRepository),
		fx.Provide(loans.NewLoanRepository),
		fx.Provide(funding.NewLoanFundingRepository),

		fx.Provide(services.NewUsersService),
		fx.Provide(services.NewLoansService),

		fx.Provide(delivery.NewLoanController),
		fx.Provide(delivery.NewUsersController),
		fx.Provide(delivery.NewAssessmentController),
		fx.Provide(delivery.NewLoanFundingController),

		fx.Invoke(initRoutes),
	)

	app.App = container
	return app, nil
}

// NewDatabase ...
func NewDatabase(lifecycle fx.Lifecycle, cfg *config.AppConfig) *postgre.DbEngine {
	db := postgre.NewDbService(*cfg)

	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return nil
		},
	})

	return db
}

// NewHTTPServer ...
func NewHTTPServer(
	lifecycle fx.Lifecycle,
	cfg *config.AppConfig,
) *delivery.HTTPEngine {
	httpServer := delivery.NewHTTPServer(cfg)

	httpServer.InitMiddleware()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go httpServer.Serve()
			return nil
		},
		OnStop: func(_ context.Context) error {
			return nil
		},
	})

	return httpServer
}

// Start ...
func (w *WebApplication) Start(ctx context.Context) error {
	return w.App.Start(ctx)
}

// Stop perform gracefull stop
func (w *WebApplication) Stop(ctx context.Context) error {
	return w.App.Stop(ctx)
}
