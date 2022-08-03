package event

const (
	ApplyForCreatingRoom = "ApplyForCreatingRoom" // 申请房间
	ApplyForJoiningRoom  = "ApplyForJoiningRoom"  // 申请加入房间
	ApplyPlaceThePiece   = "ApplyPlaceThePiece"   // 在棋盘上放置棋子
)

// CreateRoomEvent 创建房间事件
type CreateRoomEvent struct {
	PlayerID  string // 玩家ID
	RoomTitle string // 房间名称
}

// JoinRoomEvent 加入房间事件
type JoinRoomEvent struct {
	PlayerID string // 玩家ID
	RoomCode string // 房间短码
}

// PlaceThePieceEvent 在棋盘上落子事件
type PlaceThePieceEvent struct {
	PlayerID string // 玩家ID
	X        int    // 横坐标
	Y        int    // 纵坐标
}
