package goband

import "github.com/gorilla/websocket"

type PlayerClient struct {
	conn *websocket.Conn
}
