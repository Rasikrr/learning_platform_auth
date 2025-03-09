package app

import (
	"context"
	"github.com/Rasikrr/learning_platform_core/application"
)

type App struct {
	*application.App
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

func (a *App) Init(_ context.Context) error {
	return nil
}
