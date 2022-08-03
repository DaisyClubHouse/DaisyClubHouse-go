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

type PlaceThePieceRequest struct {
	X int `json:"x"` // 横坐标
	Y int `json:"y"` // 纵坐标
}

type BroadcastGameBeginning struct {
	RoomID      string `json:"room_id"`      // 房间ID
	WhiteHolder string `json:"white_holder"` // 白棋所有者ID
	BlackHolder string `json:"black_holder"` // 黑棋所有者ID
	WhoseTurn   string `json:"whose_turn"`   // 轮到谁下
}
