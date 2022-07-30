package room

import (
	"DaisyClubHouse/goband/player"
	"DaisyClubHouse/utils"
)

type Room struct {
	ID     string
	Owner  *player.Client
	Player *player.Client
}

func NewRoom(owner *player.Client) *Room {
	return &Room{
		ID:     utils.GenerateRandomID(),
		Owner:  owner,
		Player: nil,
	}
}
