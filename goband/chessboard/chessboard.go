package chessboard

import (
	"log"
	"sync"

	"DaisyClubHouse/goband/event"
	"DaisyClubHouse/goband/player"
	"DaisyClubHouse/goband/room"
	"github.com/asaskevich/EventBus"
)

type ChessBoard struct {
	clients map[string]*player.Client
	lock    sync.RWMutex
	Bus     EventBus.Bus
	rooms   map[string]*room.Room
}

func NewChessBoard() *ChessBoard {
	bus := EventBus.New()
	chessboard := ChessBoard{
		clients: make(map[string]*player.Client),
		lock:    sync.RWMutex{},
		Bus:     bus,
		rooms:   make(map[string]*room.Room),
	}

	err := bus.Subscribe(event.ApplyForCreateRoom, chessboard.eventApplyForCreateRoom)
	if err != nil {
		log.Fatalf("Subscribe error: %v\n", err)
		return nil
	}

	return &chessboard
}

func (b *ChessBoard) eventApplyForCreateRoom(event *event.CreateRoomEvent) {
	log.Printf("eventApplyForCreateRoom: %v\n", event)
	b.lock.Lock()
	defer b.lock.Unlock()

	client := b.clients[event.PlayerID]

	newRoom := room.NewRoom(client)
	b.rooms[newRoom.ID] = newRoom
}

func (b *ChessBoard) ClientConnected(client *player.Client) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.clients[client.ID] = client
}

func (b *ChessBoard) ClientDisconnected(client *player.Client) {
	log.Println("ClientDisconnected")
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, client.ID)
}
