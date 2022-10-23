package main

import (
	"DaisyClubHouse/gobang/manager"
	"DaisyClubHouse/infrastructure/server"
	"DaisyClubHouse/infrastructure/server/middleware"
	"DaisyClubHouse/infrastructure/ws"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// game
		fx.Provide(manager.NewGameManagerInstance),
		// server (http + ws)
		fx.Provide(ws.WebsocketAdaptor),
		fx.Provide(middleware.Authorization),
		fx.Invoke(server.HttpServer),
	)
	app.Run()
}
