package msg

type CreateRoomRequest struct {
	RoomTitle string `json:"room_title"` // 房间名称
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"` // 房间ID
}

type JoinRoomRequest struct {
	PlayerID   string `json:"player_id"`   // 玩家ID
	PlayerName string `json:"player_name"` // 玩家姓名
	RoomID     string `json:"room_id"`     // 房间ID
}

type JoinRoomResponse struct {
	RoomID string `json:"room_id"` // 房间ID
}

type PlaceThePieceRequest struct {
	X int `json:"x"` // 横坐标
	Y int `json:"y"` // 纵坐标
}

type PlayerProfile struct {
	ID     string `json:"id"`     // 玩家ID
	Name   string `json:"name"`   // 玩家姓名
	Avatar string `json:"avatar"` // 头像
}

type BroadcastGameBeginning struct {
	RoomID      string         `json:"room_id"`      // 房间ID
	WhiteHolder *PlayerProfile `json:"white_holder"` // 白棋所有者ID
	BlackHolder *PlayerProfile `json:"black_holder"` // 黑棋所有者ID
}

type BroadcastPlayerAction struct {
	RoomId    string         `json:"room_id"`    // 房间ID
	WhoseTurn *PlayerProfile `json:"whose_turn"` // 操作玩家
}

type BroadcastPlayerPlaceThePiece struct {
	RoomID     string `json:"room_id"`     // 房间ID
	PlayerID   string `json:"player_id"`   // 玩家ID
	PieceWhite bool   `json:"piece_white"` // 是否是白棋
	X          int    `json:"x"`           // 横坐标
	Y          int    `json:"y"`           // 纵坐标
}
