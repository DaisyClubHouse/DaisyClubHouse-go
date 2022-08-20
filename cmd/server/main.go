package main

import (
	"log"
	"net/http"

	"DaisyClubHouse/goband/game/gamemanaer"
	"DaisyClubHouse/goband/game/player"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, game *gamemanaer.GameManager) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}

	log.Printf("%s connected\n", c.RemoteAddr())

	// 初始化玩家客户端
	client := player.NewPlayerClient(c, game.Bus)
	game.ClientConnected(client)

	client.Run()
}

func main() {
	const address = "0.0.0.0:8000"

	game := gamemanaer.NewGameManager()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, game)
	})

	log.Printf("listening on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
