package hub

import (
	"DaisyClubHouse/gobang/client"
	"DaisyClubHouse/gobang/player"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
)

// GameHub 游戏总线
type GameHub struct {
	players     map[string]*player.Player // 已认证玩家
	connections map[*client.Client]bool   // 未认证长链接

	Register chan *websocket.Conn
	bus      EventBus.Bus // 事件总线
}

func NewGameHub(bus EventBus.Bus) *GameHub {
	return &GameHub{
		players:     map[string]*player.Player{},
		connections: map[*client.Client]bool{},
		Register:    make(chan *websocket.Conn, 10),
		bus:         bus,
	}
}

func (hub *GameHub) Certificate() {

}

func (hub *GameHub) run() {
	for {
		select {
		case conn := <-hub.Register:
			// 默认未认证
			c := client.NewClient(conn, hub.bus)
			c.Run()

			hub.connections[c] = false
		}
	}
}
