package main

import (
	"log"
	"net/http"

	"DaisyClubHouse/goband/player"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}

	log.Printf("%s connected\n", c.RemoteAddr())

	// 初始化玩家客户端
	player.NewPlayerClient(c).Run()
}

func main() {
	const address = "127.0.0.1:9000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r)
	})

	log.Printf("listening on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
