package main

import (
	"DaisyClubHouse/gobang/hub"
	"DaisyClubHouse/gobang/manager/client"
	"DaisyClubHouse/gobang/manager/game"
	"DaisyClubHouse/infrastructure/server"
	"DaisyClubHouse/infrastructure/server/handler"
	"DaisyClubHouse/infrastructure/ws"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// game hub
		fx.Provide(hub.GameHubProvider),
		// game
		fx.Provide(handler.NewHttpServerHandler),
		fx.Provide(game.NewGameManagerInstance),
		fx.Provide(client.PlayerClientManagerProvider),
		// server (http + ws)
		fx.Provide(ws.WebsocketAdaptor),
		fx.Invoke(server.HttpServer),
	)
	app.Run()
}
