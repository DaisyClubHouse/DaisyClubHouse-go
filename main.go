package main

import (
	"DaisyClubHouse/domain/aggregate"
	"DaisyClubHouse/infrastructure/server"
	"DaisyClubHouse/infrastructure/server/handler"
	"DaisyClubHouse/infrastructure/ws"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// game
		fx.Provide(aggregate.NewGameManager),
		// server (http + ws)
		fx.Provide(ws.WebsocketAdaptor),
		fx.Provide(handler.Authorization),
		fx.Invoke(server.HttpServer),
	)
	app.Run()
}
