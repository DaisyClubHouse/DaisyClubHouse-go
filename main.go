package main

import (
	"DaisyClubHouse/gobang/hub"
	"DaisyClubHouse/gobang/manager/client"
	"DaisyClubHouse/gobang/manager/game"
	"DaisyClubHouse/infrastructure/bus"
	"DaisyClubHouse/infrastructure/server"
	"DaisyClubHouse/infrastructure/server/handler"
	"DaisyClubHouse/infrastructure/ws"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// 事件总线
		fx.Provide(bus.EventBusProvider),
		// game hub
		fx.Provide(hub.NewGameHub),
		fx.Invoke(hub.GameHubInvoker),
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
