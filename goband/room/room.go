package room

import (
	"log"
	"sync"

	"DaisyClubHouse/goband/msg"
	"DaisyClubHouse/goband/player"
	"DaisyClubHouse/utils"
)

type Room struct {
	ID          string
	Owner       *player.Client
	Player      *player.Client
	lock        sync.RWMutex
	status      Status
	whoseTurn   *player.Client
	whiteHolder *player.Client // 执白棋的玩家（先行）
	blackHolder *player.Client // 执黑棋的玩家（后行）
}

type Status int

const (
	Status_Waiting  = iota // 等待玩家加入
	Status_Playing         // 游戏中
	Status_Finished        // 游戏结束
)

func NewRoom(owner *player.Client) *Room {
	room := &Room{
		ID:          utils.GenerateRandomID(),
		Owner:       owner,
		Player:      nil,
		lock:        sync.RWMutex{},
		status:      Status_Waiting,
		whoseTurn:   nil,
		whiteHolder: nil,
		blackHolder: nil,
	}

	log.Printf("NewRoom[%s] Created by Player(%s::%v)\n",
		room.ID, room.Owner.ID, room.Owner.RemoteAddr())
	return room
}

func (room *Room) PlayerJoin(player *player.Client) {
	room.lock.Lock()
	defer room.lock.Unlock()

	room.Player = player
	log.Printf("Player[%s] Joined Room[%s]\n", player.ID, room.ID)

	// 状态流转
	room.gameBegin()
}

func (room *Room) OwnerHold() string {
	if room.whiteHolder == room.Owner {
		return "执白"
	} else if room.blackHolder == room.Owner {
		return "执黑"
	} else {
		return "未知"
	}
}

func (room *Room) PlayerHold() string {
	if room.whiteHolder == room.Player {
		return "执白"
	} else if room.blackHolder == room.Player {
		return "执黑"
	} else {
		return "未知"
	}
}

func (room *Room) gameBegin() {
	room.status = Status_Playing

	// 准备开始游戏
	log.Println("----------Ready to start game----------")

	// 随机决定执黑执白
	if utils.RandomHalfRate() {
		room.whiteHolder = room.Player
		room.blackHolder = room.Owner
	} else {
		room.whiteHolder = room.Owner
		room.blackHolder = room.Player
	}

	log.Printf("|Room[%s]", room.ID)
	log.Printf("|Owner%s[%s,%s]",
		room.OwnerHold(),
		room.Owner.ID, room.Owner.RemoteAddr())
	log.Printf("|Player%s[%s,%s]",
		room.PlayerHold(),
		room.Player.ID, room.Player.RemoteAddr())

	// 通知执白执黑已经谁先行
	room.whoseTurn = room.whiteHolder

	pack := msg.UserPack[msg.BroadcastGameBeginning]{
		Type: msg.KindBroadcastRoomGameBeginning,
		Payload: msg.BroadcastGameBeginning{
			RoomID:      room.ID,
			WhiteHolder: room.whiteHolder.ID,
			BlackHolder: room.blackHolder.ID,
			WhoseTurn:   room.whiteHolder.ID,
		},
	}
	room.Broadcast(pack.Marshal())
}

func (room *Room) turnChanged() {
	room.lock.Lock()
	defer room.lock.Unlock()

	if room.whoseTurn == room.whiteHolder {
		room.whoseTurn = room.blackHolder
	} else {
		room.whoseTurn = room.whiteHolder
	}
}

// Broadcast 房间内广播
func (room *Room) Broadcast(data []byte) {
	if len(data) == 0 {
		// 如果没有数据，则不发送
		log.Println("------------------ERROR------------------")
		log.Printf("[ERROR] Room[%s] Broadcast data is empty\n", room.ID)
		return
	}

	room.lock.RLock()
	defer room.lock.RUnlock()

	if room.Player != nil {
		room.Player.SendRawMessage(data)
	}
	if room.Owner != nil {
		room.Owner.SendRawMessage(data)
	}
}
