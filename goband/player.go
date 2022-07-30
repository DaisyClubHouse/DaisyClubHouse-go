package goband

import (
	"log"

	"github.com/gorilla/websocket"
)

type PlayerClient struct {
	conn *websocket.Conn

	send chan []byte
	cb   *ChessBoard
}

func NewPlayerClient(conn *websocket.Conn, cb *ChessBoard) *PlayerClient {
	return &PlayerClient{
		conn: conn,
		send: make(chan []byte),
		cb:   cb,
	}
}

func (client *PlayerClient) ReadPump() {
	defer func() {
		client.cb.ClientDisconnected(client)
		client.conn.Close()
	}()
	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("read message error: %v", err)
			break
		}

		log.Printf("[recv from %s]: (t: %d, lens:%d) - %s",
			client.conn.RemoteAddr(), mt, len(message), message)
	}
}

func (client *PlayerClient) WritePump() {

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

func (client *PlayerClient) SendRawMessage(raw []byte) {
	client.send <- raw
}
