package ws

import (
	"log"
	"net/http"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/manager"
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

func wsHandler(w http.ResponseWriter, r *http.Request, game *manager.GameManager) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// 初始化玩家长链接
	client := entity.GeneratePlayerClient(ws, game.Bus)
	game.Connect(client)
	client.Run()

	/*
		for {
			// 读取ws中的数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("[message] New message from %s {type: %d, msg: %s}", ws.RemoteAddr(), mt, string(message))
			if string(message) == "ping" {
				message = []byte("pong")
			}
		}
	*/
}

func WebsocketAdaptor(game *manager.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, game)
	}
}
