package chessboard

import (
	"log"
	"sync"

	"DaisyClubHouse/goband/player"
)

type ChessBoard struct {
	clients []*player.Client
	lock    sync.RWMutex
}

func NewChessBoard() *ChessBoard {
	return &ChessBoard{
		clients: make([]*player.Client, 0),
		lock:    sync.RWMutex{},
	}
}

func (b *ChessBoard) ClientConnected(client *player.Client) {
	b.lock.RLock()
	b.clients = append(b.clients, client)
	defer b.lock.RUnlock()

}

func (b *ChessBoard) ClientDisconnected(client *player.Client) {
	log.Println("ClientDisconnected")
	for i, c := range b.clients {
		if c == client {
			log.Printf("客户端[%s]断开连接\n", client.RemoteAddr())
			b.clients = append(b.clients[:i], b.clients[i+1:]...)
			break
		}
	}
}
