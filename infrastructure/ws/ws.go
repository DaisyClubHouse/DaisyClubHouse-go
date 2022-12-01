package ws

import (
	"log"
	"net/http"

	"DaisyClubHouse/gobang/hub"
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

func wsHandler(w http.ResponseWriter, r *http.Request, hub *hub.GameHub) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("upgrade error: %v", err)
		return
	}

	// 建立链接
	hub.Connect(ws)
}

func WebsocketAdaptor(hub *hub.GameHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, hub)
	}
}
