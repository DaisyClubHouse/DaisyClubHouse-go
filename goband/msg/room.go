package msg

type CreateRoomRequest struct {
	RoomTitle string `json:"room_title"` // 房间名称
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"` // 房间ID
}

type JoinRoomRequest struct {
	ShortCode string `json:"short_code"` // 房间短码
}

type JoinRoomResponse struct {
	RoomID string `json:"room_id"` // 房间ID
}
