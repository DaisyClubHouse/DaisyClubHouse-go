package main

import (
	"DaisyClubHouse/goband/game/gamemanaer"
	"DaisyClubHouse/infrastructure/server"
	"DaisyClubHouse/infrastructure/ws"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// game
		fx.Provide(gamemanaer.NewGameManager),
		// server (http + ws)
		fx.Provide(ws.WebsocketAdaptor),
		fx.Invoke(server.HttpServer),
	)
	app.Run()
}
