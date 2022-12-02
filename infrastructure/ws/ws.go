package ws

import (
	"log"
	"net/http"

	"DaisyClubHouse/gobang/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketAdaptor(hub *hub.GameHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatalf("upgrade error: %v", err)
			return
		}

		slog.Info("receive a new connection", slog.Any("remote_addr", ws.RemoteAddr()))

		// 建立链接
		_ = hub.Connect(ws)
	}
}
