package goband

type Room struct {
	ID     string
	Owner  *PlayerClient
	Player *PlayerClient
}
