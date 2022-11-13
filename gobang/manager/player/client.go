package player

import (
	"log"
	"net"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/msg"
	"DaisyClubHouse/utils"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

type Client struct {
	ID   string
	conn *websocket.Conn

	send     chan []byte
	close    chan struct{}
	bus      EventBus.Bus
	Identity *entity.UserInfo // 认证身份信息
}

func GeneratePlayerClient(conn *websocket.Conn, bus EventBus.Bus) *Client {
	return &Client{
		ID:    utils.GenerateRandomID(),
		conn:  conn,
		send:  make(chan []byte),
		close: make(chan struct{}),
		bus:   bus,
	}
}

func (client *Client) Run() {
	client.conn.SetCloseHandler(func(code int, text string) error {
		slog.Info("客户端断开连接",
			slog.String("client_id", client.ID),
			slog.Any("addr", client.conn.RemoteAddr()),
			slog.Int("code", code),
			slog.String("text", text))

		// 关闭通道
		close(client.close)
		close(client.send)
		return nil
	})

	slog.Info("客户端已连接", slog.String("client_id", client.ID), slog.Any("addr", client.conn.RemoteAddr()))

	go client.writePumpLoop()
	client.readPumpLoop()
}

func (client *Client) RemoteAddr() net.Addr {
	return client.conn.RemoteAddr()
}

func (client *Client) PlayerID() string {
	return client.Identity.ID
}

func (client *Client) PlayerProfile() *msg.PlayerProfile {
	return &msg.PlayerProfile{
		ID:     client.Identity.ID,
		Name:   client.Identity.Username,
		Avatar: client.Identity.Avatar,
	}
}

func (client *Client) readPumpLoop() {
	defer func() {
		client.conn.Close()
	}()

	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("链接已断开", slog.String("client_id", client.ID))
				break
			}
			slog.Error("读取消息失败", err, slog.String("client_id", client.ID))
			break
		}

		slog.Info("receive message",
			slog.String("client_id", client.ID),
			slog.Any("addr", client.conn.RemoteAddr()),
			slog.Int("mType", mt),
			slog.String("msg", string(message)),
		)

		if mt == websocket.TextMessage {
			kind, payload, err := msg.Parsing(message)
			if err != nil {
				slog.Error("消息解析失败", err)
				return
			}

			handlerMap := eventHandleDispatcher()
			if handler, ok := handlerMap[kind]; ok {
				handler(client, payload)
			} else {
				slog.Warn("Unknown message, DISCARD!",
					slog.String("client_id", client.ID),
					slog.Any("addr", client.conn.RemoteAddr()),
					slog.Any("kind", kind),
					slog.String("payload", string(payload)),
				)
			}
		}
	}
}

func (client *Client) writePumpLoop() {
	defer func() {
		client.conn.Close()
	}()

	heartbeatTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				slog.Warn("链接已关闭", slog.String("client_id", client.ID), slog.Any("addr", client.conn.RemoteAddr()))
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
			// log.Printf("[send] Ping %v", client.conn.RemoteAddr())
			if err := client.conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
				log.Printf("[error] WriteMessage: %v", err)
				return
			}
		}
	}
}

func (client *Client) SendRawMessage(raw []byte) {
	log.Printf("[send to %s] address:%s, lens:%d", client.ID, client.RemoteAddr(), len(raw))
	client.send <- raw
}
