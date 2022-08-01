package player

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"DaisyClubHouse/goband/event"
	"DaisyClubHouse/goband/msg"
	"DaisyClubHouse/utils"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	conn *websocket.Conn

	send  chan []byte
	close chan struct{}
	bus   EventBus.Bus
}

func NewPlayerClient(conn *websocket.Conn, bus EventBus.Bus) *Client {
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
		log.Printf("[close] 客户端[%v]断开了连接(code:%d,text:%s)",
			client.conn.RemoteAddr(), code, text)
		// 关闭通道
		close(client.close)
		close(client.send)
		return nil
	})

	go client.writePump()
	go client.readPump()
}

func (client *Client) RemoteAddr() net.Addr {
	return client.conn.RemoteAddr()
}

func (client *Client) readPump() {
	defer func() {
		client.conn.Close()
	}()

	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[error] Client Closed: %v", err)
				break
			}
			log.Printf("read message error: %v", err)
			break
		}

		log.Printf("[recv from %s]: (t: %d, lens:%d) - %s",
			client.conn.RemoteAddr(), mt, len(message), message)

		if mt == websocket.TextMessage {
			kind, payload, err := msg.Parsing(message)
			if err != nil {
				return
			}
			switch kind {
			case msg.KindCreateRoomRequest:
				var req msg.CreateRoomRequest
				if err := json.Unmarshal(payload, &req); err != nil {
					log.Printf("[error] json.Unmarshal: %v", err)
					return
				}

				evt := event.CreateRoomEvent{
					PlayerID:  client.ID,
					RoomTitle: req.RoomTitle,
				}
				client.bus.Publish(event.ApplyForCreatingRoom, &evt)
			case msg.KindJoinRoomRequest:
				var req msg.JoinRoomRequest
				if err := json.Unmarshal(payload, &req); err != nil {
					log.Printf("[error] json.Unmarshal: %v", err)
					return
				}

				evt := event.JoinRoomEvent{
					PlayerID: client.ID,
					RoomCode: req.ShortCode,
				}
				client.bus.Publish(event.ApplyForJoiningRoom, &evt)
			}
		}
	}
}

func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()

	heartbeatTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				log.Println("[error] 发送信息出错，关闭连接")
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("[error] GetWriter: %v", err)
				return
			}
			w.Write(message)

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

func (client *Client) sendRawMessage(raw []byte) {
	client.send <- raw
}

// ApplyForCreatingRoom 发送创建房间请求
// func (client *Client) ApplyForCreatingRoom() error {
// 	req := pb.CreateRoomRequest{
// 		RoomTitle: "test1",
// 	}
// 	bytes, err := proto.Marshal(&req)
// 	if err != nil {
// 		return err
// 	}
//
// 	msgPack := pb.UserMsgPack{
// 		Kind: pb.MsgKind_MsgCreateRoomRequest,
// 		Data: bytes,
// 	}
//
// 	raw, err := proto.Marshal(&msgPack)
// 	if err != nil {
// 		return err
// 	}
//
// 	client.sendRawMessage(raw)
// 	return nil
// }
//
// // ApplyForJoiningRoom 发送加入房间请求
// func (client *Client) ApplyForJoiningRoom() error {
//
// }
