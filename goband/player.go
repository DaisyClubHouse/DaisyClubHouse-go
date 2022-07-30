package goband

import (
	"log"

	"github.com/gorilla/websocket"
)

type PlayerClient struct {
	conn *websocket.Conn
}

func NewPlayerClient(conn *websocket.Conn) *PlayerClient {
	return &PlayerClient{conn: conn}
}

func (client *PlayerClient) run() {
	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Fatalf("read error: %v", err)
			break
		}

		log.Printf("recv: %s", message)
		if err = client.conn.WriteMessage(mt, message); err != nil {
			log.Fatalf("write error: %v", err)
			break
		}
	}
}
