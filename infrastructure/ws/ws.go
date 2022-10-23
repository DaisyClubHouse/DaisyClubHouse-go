package ws

import (
	"log"
	"net/http"

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

	log.Printf("%s connected\n", ws.RemoteAddr())
	for {
		// 读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		// 写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}

	// 初始化玩家客户端
	// client := player.NewPlayerClient(ws, game.Bus)
	// game.ClientConnected(client)
	//
	// client.Run()
}

func WebsocketAdaptor(game *manager.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, game)
	}
}
