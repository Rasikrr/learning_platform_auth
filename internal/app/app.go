package app

import (
	"context"
	authC "github.com/Rasikrr/learning_platform_auth/internal/cache/auth"
	usersC "github.com/Rasikrr/learning_platform_auth/internal/clients/users"
	"github.com/Rasikrr/learning_platform_auth/internal/envs"
	"github.com/Rasikrr/learning_platform_auth/internal/ports/grpc"
	authS "github.com/Rasikrr/learning_platform_auth/internal/services/auth"
	"github.com/Rasikrr/learning_platform_core/application"
)

type App struct {
	*application.App
	authCache   authC.Cache
	usersClient usersC.Client
	authService authS.Service
}

func NewApp(ctx context.Context, name string) (*App, error) {
	app := &App{
		App: application.NewApp(ctx, name),
	}
	if err := app.Init(ctx); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Init(ctx context.Context) error {
	for _, init := range []func(context.Context) error{
		a.initCache,
		a.initClients,
		a.initServices,
		a.initGRPCServer,
	} {
		if err := init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initCache(_ context.Context) error {
	a.authCache = authC.NewCache(a.Redis())
	return nil
}

func (a *App) initClients(ctx context.Context) error {
	var err error
	a.usersClient, err = usersC.NewClient(ctx, a.Config().Env.Get(envs.UsersGRPcAddress).GetString())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServices(_ context.Context) error {
	a.authService = authS.NewService(
		a.Config().Env.Get(envs.AccessTokenTTL).GetDuration(),
		a.Config().Env.Get(envs.RefreshTokenTTL).GetDuration(),
		a.authCache,
		a.usersClient,
	)
	return nil
}

func (a *App) initGRPCServer(_ context.Context) error {
	grpc.NewServer(a.GrpcServer().Srv(), a.authService)
	return nil
}
