package player

import (
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn

	send  chan []byte
	close chan struct{}
}

func NewPlayerClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:  conn,
		send:  make(chan []byte),
		close: make(chan struct{}),
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

	go client.WritePump()
	go client.ReadPump()
}

func (client *Client) RemoteAddr() net.Addr {
	return client.conn.RemoteAddr()
}

func (client *Client) ReadPump() {
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
	}
}

func (client *Client) WritePump() {
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

func (client *Client) SendRawMessage(raw []byte) {
	client.send <- raw
}
