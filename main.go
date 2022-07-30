package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("read error: %v", err)
			break
		}

		log.Printf("recv: %s", message)
		if err = c.WriteMessage(mt, message); err != nil {
			log.Fatalf("write error: %v", err)
			break
		}
	}
}
func main() {
	const address = "127.0.0.1:9000"
	http.HandleFunc("/", handler)

	log.Printf("listening on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
