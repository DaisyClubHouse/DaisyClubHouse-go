package chessboard

import (
	"log"
	"sync"

	"DaisyClubHouse/goband/event"
	"DaisyClubHouse/goband/game/player"
	"DaisyClubHouse/goband/game/room"
	"DaisyClubHouse/utils"
	"github.com/asaskevich/EventBus"
)

type ChessBoard struct {
	clients         map[string]*player.Client
	lock            sync.RWMutex
	Bus             EventBus.Bus
	rooms           map[string]*room.Room
	codeRoomMapping map[string]string // code -> roomID
}

func NewChessBoard() *ChessBoard {
	bus := EventBus.New()
	chessboard := ChessBoard{
		clients:         make(map[string]*player.Client),
		lock:            sync.RWMutex{},
		Bus:             bus,
		rooms:           make(map[string]*room.Room),
		codeRoomMapping: make(map[string]string),
	}

	err := bus.Subscribe(event.ApplyForCreatingRoom, chessboard.eventApplyForCreatingRoom)
	if err != nil {
		log.Fatalf("Subscribe error: %v\n", err)
		return nil
	}

	err = bus.Subscribe(event.ApplyForJoiningRoom, chessboard.eventApplyForJoiningRoom)
	if err != nil {
		log.Fatalf("Subscribe error: %v\n", err)
		return nil
	}

	return &chessboard
}

// 创建房间处理事件
func (b *ChessBoard) eventApplyForCreatingRoom(event *event.CreateRoomEvent) {
	log.Printf("eventApplyForCreatingRoom: %v\n", event)
	b.lock.Lock()
	defer b.lock.Unlock()

	client := b.clients[event.PlayerID]

	newRoom := room.NewRoom(client)
	b.rooms[newRoom.ID] = newRoom

	// 生成随机code映射
	code := utils.GenerateSixFigure()
	b.codeRoomMapping[code] = newRoom.ID

	log.Printf("【创建新房间】code: %s, roomID: %s\n", code, newRoom.ID)
}

// 加入房间处理事件
func (b *ChessBoard) eventApplyForJoiningRoom(e *event.JoinRoomEvent) {
	log.Printf("eventApplyForJoiningRoom: %v\n", e)
	b.lock.Lock()
	defer b.lock.Unlock()

	roomId, ok := b.codeRoomMapping[e.RoomCode]
	if !ok {
		// TODO 通知玩家未找到房间
		log.Printf("通知玩家未找到房间")
		log.Println(b.codeRoomMapping)
		return
	}

	client := b.clients[e.PlayerID]
	targetRoom := b.rooms[roomId]
	// 玩家加入
	targetRoom.PlayerJoin(client)
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
