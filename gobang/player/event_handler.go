package player

import (
	"encoding/json"

	"DaisyClubHouse/gobang/message/receiver"
	"DaisyClubHouse/gobang/message/user_pack"
	"golang.org/x/exp/slog"
)

func eventHandleDispatcher() map[user_pack.Kind]func(*Player, []byte) {
	handlers := map[user_pack.Kind]func(*Player, []byte){
		user_pack.KindJoinRoomRequest:      handleJoinRoomRequest,
		user_pack.KindPlaceThePieceRequest: handlePlaceThePieceRequest,
	}

	return handlers
}

// handleJoinRoomRequest 处理进入房间请求
func handleJoinRoomRequest(client *Player, payload []byte) {
	var req user_pack.JoinRoomRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("消息反序列化失败", err, slog.String("payload", string(payload)))
		return
	}

	evt := receiver.JoinRoomEvent{
		PlayerID:   req.PlayerID,
		PlayerName: req.PlayerName,
		ClientID:   client.ID,
		RoomID:     req.RoomID,
	}
	client.bus.Publish(receiver.ApplyForJoiningRoom, &evt)
}

// handlePlaceThePieceRequest 处理五子棋落子请求
func handlePlaceThePieceRequest(client *Player, payload []byte) {
	var req user_pack.PlaceThePieceRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("消息反序列化失败", err, slog.String("payload", string(payload)))
		return
	}
	evt := receiver.PlaceThePieceEvent{
		ClientID: client.ID,
		X:        req.X,
		Y:        req.Y,
	}
	client.bus.Publish(receiver.ApplyPlaceThePiece, &evt)
}
