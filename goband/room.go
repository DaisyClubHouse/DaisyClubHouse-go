package goband

import "DaisyClubHouse/goband/player"

type Room struct {
	ID     string
	Owner  *player.Client
	Player *player.Client
}
