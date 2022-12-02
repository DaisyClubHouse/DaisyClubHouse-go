package client

import (
	"log"
	"time"

	"DaisyClubHouse/gobang/message/user_pack"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

const (
	writeWait            = 10 * time.Second
	heartBeatPeriod      = 10 * time.Second // 心跳轮训时间间隔
	maxReadMessageLength = 512              // 最大可读消息长度 bytes
)

type Client struct {
	conn *websocket.Conn
	send chan []byte

	bus EventBus.Bus // 事件总线
}

func NewClient(conn *websocket.Conn, bus EventBus.Bus) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte),
		bus:  bus,
	}

}

func (client *Client) Run() {
	go client.readPumpLoop()
	go client.writePumpLoop()
}

// readPumpLoop 读循环
func (client *Client) readPumpLoop() {
	client.conn.SetReadLimit(maxReadMessageLength)

	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("链接已断开", slog.Any("client_id", client.conn.RemoteAddr()))
				break
			}
			slog.Error("读取消息失败", err, slog.Any("remote_addr", client.conn.RemoteAddr()))
			break
		}

		slog.Info("receive message",
			// slog.String("client_id", client.ID),
			slog.Any("addr", client.conn.RemoteAddr()),
			slog.Int("mType", mt),
			slog.String("msg", string(message)),
		)

		if mt == websocket.TextMessage {
			kind, payload, err := user_pack.Parsing(message)
			if err != nil {
				slog.Error("消息解析失败", err)
				return
			}

			// handlerMap := eventHandleDispatcher()
			// if handler, ok := handlerMap[kind]; ok {
			// 	handler(client, payload)
			// } else {
			slog.Warn("Unknown message, DISCARD!",
				slog.Any("addr", client.conn.RemoteAddr()),
				slog.Any("kind", kind),
				slog.String("payload", string(payload)),
			)
			// }
		}
	}
}

// writePumpLoop 写循环
func (client *Client) writePumpLoop() {
	heartbeatTicker := time.NewTicker(heartBeatPeriod)

	defer func() {
		heartbeatTicker.Stop()
		client.conn.Close()

		slog.Info("writePumpLoop exit",
			slog.Any("addr", client.conn.RemoteAddr()))
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				slog.Warn("链接已关闭", slog.Any("addr", client.conn.RemoteAddr()))
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[error] GetWriter: %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("[error] Write message: %v", err)
				return
			}

			if err := w.Close(); err != nil {
				log.Printf("[error] Close Writer: %v", err)
				return
			}
		case <-heartbeatTicker.C:
			log.Printf("[send] Ping %v", client.conn.RemoteAddr())
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[error] WriteMessage: %v", err)
				return
			}
		}
	}

}
