package event

const (
	ApplyForCreatingRoom = "ApplyForCreatingRoom" // 申请房间
	ApplyForJoiningRoom  = "ApplyForJoiningRoom"  // 申请加入房间
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
