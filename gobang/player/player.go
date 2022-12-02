package player

import (
	"log"
	"net"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/message/inner"
	"DaisyClubHouse/gobang/message/receiver"
	"DaisyClubHouse/gobang/message/user_pack"
	"DaisyClubHouse/utils"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

// Player 玩家对象
type Player struct {
	ID       string
	UserInfo *entity.UserInfo // 认证身份信息
	online   bool             // 长链接是否在线
	conn     *websocket.Conn
	send     chan []byte
	close    chan struct{}
	bus      EventBus.Bus
}

func Wrap2Player(conn *websocket.Conn) *Player {
	return &Player{
		ID:       utils.GenerateRandomID(),
		UserInfo: nil,
		online:   true,
		conn:     conn,
		send:     make(chan []byte, 256),
		close:    make(chan struct{}),
		// bus:      bus,
	}
}

func (client *Player) Run() {
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
	go client.readPumpLoop()
}

func (client *Player) RemoteAddr() net.Addr {
	return client.conn.RemoteAddr()
}

func (client *Player) PlayerID() string {
	return client.UserInfo.ID
}

func (client *Player) PlayerProfile() *user_pack.PlayerProfile {
	return &user_pack.PlayerProfile{
		ID:     client.UserInfo.ID,
		Name:   client.UserInfo.Username,
		Avatar: client.UserInfo.Avatar,
	}
}

func (client *Player) Disconnect() {
	client.NoticedDisconnect()

	// 房间断线通知
	client.bus.Publish(receiver.PlayerDisconnect, &inner.PlayerDisconnectEvent{ClientID: client.ID})
}

func (client *Player) NoticedDisconnect() {
	if !client.online {
		return
	}

	slog.Info("关闭链接", slog.String("player_id", client.PlayerID()), slog.String("client_id", client.ID))

	if err := client.conn.Close(); err != nil {
		slog.Error("链接关闭错误", err)
	}

	client.close <- struct{}{}
}

// SendRawMessage 给玩家推送消息
func (client *Player) SendRawMessage(raw []byte) {
	slog.Info("[send msg to player]",
		slog.String("id", client.ID),
		slog.Any("remoteAddr", client.RemoteAddr()),
		slog.Int("length", len(raw)))

	client.send <- raw
}

// readPumpLoop 读循环
func (client *Player) readPumpLoop() {
	defer func() {
		client.Disconnect()
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
			kind, payload, err := user_pack.Parsing(message)
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

// writePumpLoop 写循环
func (client *Player) writePumpLoop() {
	defer func() {
		slog.Info("writePumpLoop exit",
			slog.String("player_id", client.PlayerID()),
			slog.String("client_id", client.ID),
			slog.Any("addr", client.conn.RemoteAddr()))

		client.Disconnect()
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
		case <-client.close:
			client.online = false
			return
		case <-heartbeatTicker.C:
			log.Printf("[send] Ping %v", client.conn.RemoteAddr())
			if err := client.conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
				log.Printf("[error] WriteMessage: %v", err)
				return
			}
		}
	}

}
