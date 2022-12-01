package hub

import (
	"DaisyClubHouse/gobang/player"
	"github.com/gorilla/websocket"
)

// GameHub 游戏总线
type GameHub struct {
	players map[string]*player.Player // 已认证玩家

	register   chan *player.Player
	unregister chan *player.Player
}

func NewGameHub() *GameHub {
	return &GameHub{
		players:    map[string]*player.Player{},
		register:   make(chan *player.Player),
		unregister: make(chan *player.Player),
	}
}

func (hub *GameHub) Connect(conn *websocket.Conn) {

}

func (hub *GameHub) run() {
	for {
		select {
		case p := <-hub.register:
			// 注册
			hub.players[p.ID] = p
		case p := <-hub.unregister:
			// 登出
			delete(hub.players, p.ID)
		}
	}
}
