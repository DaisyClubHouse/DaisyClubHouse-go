package goband

import "github.com/gorilla/websocket"

type ChessBoard struct {
	clients []*PlayerClient
}

func NewChessBoard() *ChessBoard {
	return &ChessBoard{}
}

func (b *ChessBoard) ClientConnected(conn *websocket.Conn) {
	client := NewPlayerClient(conn)
	b.clients = append(b.clients, client)

	client.run()
}
