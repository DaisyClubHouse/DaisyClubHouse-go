package room

import (
	"log"
	"sync"

	"DaisyClubHouse/goband/player"
	"DaisyClubHouse/utils"
)

type Room struct {
	ID     string
	Owner  *player.Client
	Player *player.Client
	lock   sync.RWMutex
}

func NewRoom(owner *player.Client) *Room {
	room := &Room{
		ID:     utils.GenerateRandomID(),
		Owner:  owner,
		Player: nil,
		lock:   sync.RWMutex{},
	}

	log.Printf("NewRoom[%s] Created by Player(%s::%v)\n",
		room.ID, room.Owner.ID, room.Owner.RemoteAddr())
	return room
}

func (room *Room) PlayerJoin(player *player.Client) {
	room.lock.Lock()
	room.Player = player
	room.lock.Unlock()

	log.Printf("Player[%s] Joined Room[%s]\n", player.ID, room.ID)

	// TODO 准备开始游戏
	log.Println("----------Ready to start game----------")
	log.Printf("|Room[%s] Owner: %s\n, Player: %s\n", room.ID, room.Owner.ID, room.Player.ID)
	log.Println("----------Ready to start game----------")
}
