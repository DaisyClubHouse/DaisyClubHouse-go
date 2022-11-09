package player

import (
	"encoding/json"

	"DaisyClubHouse/gobang/event"
	"DaisyClubHouse/gobang/msg"
	"golang.org/x/exp/slog"
)

func eventHandleDispatcher() map[msg.Kind]func(*Client, []byte) {
	handlers := map[msg.Kind]func(*Client, []byte){
		msg.KindCreateRoomRequest:    handleCreateRoomRequest,
		msg.KindJoinRoomRequest:      handleJoinRoomRequest,
		msg.KindPlaceThePieceRequest: handlePlaceThePieceRequest,
	}

	return handlers
}

// handleCreateRoomRequest 处理创建房间请求
func handleCreateRoomRequest(client *Client, payload []byte) {
	var req msg.CreateRoomRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("消息反序列化失败", err, slog.String("payload", string(payload)))
		return
	}

	evt := event.CreateRoomEvent{
		PlayerID:  client.ID,
		RoomTitle: req.RoomTitle,
	}
	client.bus.Publish(event.ApplyForCreatingRoom, &evt)
}

// handleJoinRoomRequest 处理进入房间请求
func handleJoinRoomRequest(client *Client, payload []byte) {
	var req msg.JoinRoomRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("消息反序列化失败", err, slog.String("payload", string(payload)))
		return
	}

	evt := event.JoinRoomEvent{
		PlayerID: client.ID,
		RoomCode: req.ShortCode,
	}
	client.bus.Publish(event.ApplyForJoiningRoom, &evt)
}

// handlePlaceThePieceRequest 处理五子棋落子请求
func handlePlaceThePieceRequest(client *Client, payload []byte) {
	var req msg.PlaceThePieceRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("消息反序列化失败", err, slog.String("payload", string(payload)))
		return
	}
	evt := event.PlaceThePieceEvent{
		PlayerID: client.ID,
		X:        req.X,
		Y:        req.Y,
	}
	client.bus.Publish(event.ApplyPlaceThePiece, &evt)
}
