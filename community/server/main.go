package main

import (
	"community-server/internal/di"
	"community-server/internal/logger"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		di.Module,
		fx.Invoke(logger.InitLogger),
		fx.Invoke(di.OnStartHook),
	)

	app.Run()
}
