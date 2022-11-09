package manager

import (
	"log"
	"sync"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/event"
	"DaisyClubHouse/utils"
	"github.com/asaskevich/EventBus"
)

type GameManager struct {
	clients           map[string]*entity.Client
	lock              sync.RWMutex
	Bus               EventBus.Bus
	rooms             map[string]*entity.Room
	codeRoomMapping   map[string]string // code -> roomID
	playerRoomMapping map[string]string // playerID -> roomID
}

var once sync.Once

func NewGameManagerInstance() *GameManager {
	var gm *GameManager

	once.Do(func() {
		gm = func() *GameManager {
			bus := EventBus.New()
			chessboard := GameManager{
				clients:           make(map[string]*entity.Client),
				lock:              sync.RWMutex{},
				Bus:               bus,
				rooms:             make(map[string]*entity.Room),
				codeRoomMapping:   make(map[string]string),
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

	roomId, ok := b.codeRoomMapping[e.RoomCode]
	if !ok {
		log.Printf("通知玩家未找到房间")
		log.Println(b.codeRoomMapping)
		return
	}

	client := b.clients[e.PlayerID]
	targetRoom := b.rooms[roomId]
	// 玩家加入
	targetRoom.PlayerJoin(client)

	// 生成玩家房间映射
	b.playerRoomMapping[client.ID] = roomId
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
func (b *GameManager) Connect(client *entity.Client) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.clients[client.ID] = client
}

// Disconnect 客户端断开连接
func (b *GameManager) Disconnect(client *entity.Client) {
	log.Println("Disconnect")
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, client.ID)
}

// RoomProfileList 查询房间简要信息列表
func (b *GameManager) RoomProfileList() []entity.RoomProfile {
	profileList := make([]entity.RoomProfile, 0, len(b.rooms))
	for _, room := range b.rooms {
		profileList = append(profileList, room.RoomProfile)
	}
	return profileList
}

// CreateRoom 创建房间
func (b *GameManager) CreateRoom(user *entity.UserInfo) (*entity.RoomProfile, error) {
	log.Printf("CreateRoom: %v\n", user)
	b.lock.Lock()
	defer b.lock.Unlock()

	newRoom := entity.CreateNewRoom(user)
	b.rooms[newRoom.ID] = newRoom

	// 生成随机code映射
	code := utils.GenerateSixFigure()
	b.codeRoomMapping[code] = newRoom.ID

	log.Printf("【创建新房间】code: %s, roomID: %s\n", code, newRoom.ID)

	return &newRoom.RoomProfile, nil
}
