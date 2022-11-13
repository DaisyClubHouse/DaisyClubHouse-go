package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/event"
	"DaisyClubHouse/gobang/manager/client"
	"DaisyClubHouse/gobang/manager/player"
	"DaisyClubHouse/gobang/manager/room"
	"github.com/asaskevich/EventBus"
	"golang.org/x/exp/slog"
)

type GameManager struct {
	clientManager     *client.PlayerClientManager
	lock              sync.RWMutex
	Bus               EventBus.Bus
	rooms             map[string]*room.Room
	playerRoomMapping map[string]string // playerID -> roomID
}

var once sync.Once

func NewGameManagerInstance(clientManager *client.PlayerClientManager) *GameManager {
	var gm *GameManager

	once.Do(func() {
		gm = func() *GameManager {
			bus := EventBus.New()
			chessboard := GameManager{
				clientManager:     clientManager,
				lock:              sync.RWMutex{},
				Bus:               bus,
				rooms:             make(map[string]*room.Room),
				playerRoomMapping: make(map[string]string),
			}

			err := bus.Subscribe(event.ApplyForJoiningRoom, chessboard.eventApplyForJoiningRoom)
			if err != nil {
				log.Fatalf("Subscribe error: %v\n", err)
				return nil
			}

			err = bus.Subscribe(event.ApplyPlaceThePiece, chessboard.eventApplyPlaceThePiece)
			if err != nil {
				log.Fatalf("Subscribe error: %v\n", err)
				return nil
			}
			return &chessboard
		}()
	})

	log.Println("初始化GameManager")

	return gm
}

// 加入房间处理事件
func (b *GameManager) eventApplyForJoiningRoom(e *event.JoinRoomEvent) {
	log.Printf("eventApplyForJoiningRoom: %v\n", e)
	b.lock.Lock()
	defer b.lock.Unlock()

	targetRoom, ok := b.rooms[e.RoomID]
	if !ok {
		slog.Warn("通知玩家未找到房间", slog.String("room_id", e.RoomID))
		return
	}

	slog.Info("查找到房间",
		slog.String("room_id", targetRoom.ID),
		slog.String("title", targetRoom.Title),
	)
	// 建立playerID、clientID的关联
	b.clientManager.AssociatedID(e.ClientID, e.PlayerID)

	playerC, err := b.clientManager.GetClientByPlayerID(e.PlayerID)
	if err != nil {
		slog.Error("未找到玩家", err, slog.String("player_id", e.PlayerID))
		return
	}
	// 生成玩家房间映射
	b.playerRoomMapping[playerC.ID] = e.RoomID

	// 玩家加入房间
	targetRoom.PlayerJoin(&entity.UserInfo{
		ID:         e.PlayerID,
		Username:   e.PlayerName,
		Nickname:   "",
		Avatar:     fmt.Sprintf("https://joeschmoe.io/api/v1/%s", e.PlayerName),
		CreateTime: time.Now(),
	}, playerC)
}

// 在棋盘上落子处理事件
func (b *GameManager) eventApplyPlaceThePiece(e *event.PlaceThePieceEvent) {
	log.Printf("eventApplyPlaceThePiece: %v\n", e)
	b.lock.Lock()
	defer b.lock.Unlock()
	roomId := b.playerRoomMapping[e.PlayerID]
	targetRoom, ok := b.rooms[roomId]
	if !ok {
		log.Printf("未找到房间")
		return
	}

	result := targetRoom.PlayerPlaceThePiece(e.PlayerID, e.X, e.Y)
	log.Println(result)
}

// Connect 连接到新客户端
func (b *GameManager) Connect(client *player.Client) {
	b.clientManager.ClientConnected(client)
}

// Disconnect 客户端断开连接
func (b *GameManager) Disconnect(client *player.Client) {
	b.clientManager.ClientDisconnected(client.ID)
}

// RoomProfileList 查询房间简要信息列表
func (b *GameManager) RoomProfileList() []room.RoomProfile {
	profileList := make([]room.RoomProfile, 0, len(b.rooms))
	for _, item := range b.rooms {
		profileList = append(profileList, item.RoomProfile)
	}
	return profileList
}

// CreateRoom 创建房间
func (b *GameManager) CreateRoom(user *entity.UserInfo) (*room.RoomProfile, error) {
	log.Printf("CreateRoom: %v\n", user)
	b.lock.Lock()
	defer b.lock.Unlock()

	newRoom := room.CreateNewRoom(user)
	b.rooms[newRoom.ID] = newRoom

	log.Printf("【创建新房间】roomID: %s\n", newRoom.ID)

	return &newRoom.RoomProfile, nil
}
