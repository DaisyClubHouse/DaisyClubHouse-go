package ws

import (
	"log"
	"net/http"

	"DaisyClubHouse/gobang/manager/game"
	"DaisyClubHouse/gobang/player"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, game *game.GameManager) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// 初始化玩家长链接
	client := player.GeneratePlayerClient(ws, game.Bus)
	game.Connect(client)
	client.Run()
}

func WebsocketAdaptor(game *game.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, game)
	}
}
