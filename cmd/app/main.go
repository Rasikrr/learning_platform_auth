package main

import (
	"context"
	"github.com/Rasikrr/learning_platform_auth/internal/app"
)

const (
	appName = "auth"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx, appName)
	if err != nil {
		panic(err)
	}
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
