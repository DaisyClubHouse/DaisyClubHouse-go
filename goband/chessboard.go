package goband

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChessBoard struct {
	clients []*PlayerClient
}

func NewChessBoard() *ChessBoard {
	return &ChessBoard{}
}

func (b *ChessBoard) ClientConnected(conn *websocket.Conn) {
	client := NewPlayerClient(conn, b)
	client.WritePump()
	client.ReadPump()

	b.clients = append(b.clients, client)
}

func (b *ChessBoard) ClientDisconnected(client *PlayerClient) {
	for i, c := range b.clients {
		if c == client {
			log.Printf("客户端[%s]断开连接\n", client.conn.RemoteAddr())
			b.clients = append(b.clients[:i], b.clients[i+1:]...)
			break
		}
	}
}
