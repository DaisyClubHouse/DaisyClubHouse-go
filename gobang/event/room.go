package event

const (
	ApplyForJoiningRoom = "ApplyForJoiningRoom" // 申请加入房间
	ApplyPlaceThePiece  = "ApplyPlaceThePiece"  // 在棋盘上放置棋子
)

// JoinRoomEvent 加入房间事件
type JoinRoomEvent struct {
	PlayerID   string // 玩家ID
	PlayerName string // 玩家姓名
	ClientID   string // 客户端ID
	RoomID     string // 房间ID
}

// PlaceThePieceEvent 在棋盘上落子事件
type PlaceThePieceEvent struct {
	PlayerID string // 玩家ID
	X        int    // 横坐标
	Y        int    // 纵坐标
}
