package main

import (
	"log"
	"net/http"

	"DaisyClubHouse/goband"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(board *goband.ChessBoard, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}
	defer c.Close()

	log.Printf("%s connected\n", c.RemoteAddr())

	board.ClientConnected(c)

}
func main() {
	const address = "127.0.0.1:9000"

	chessboard := goband.NewChessBoard()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(chessboard, w, r)
	})

	log.Printf("listening on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
