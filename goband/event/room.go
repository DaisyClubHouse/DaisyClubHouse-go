package event

const (
	ApplyForCreateRoom = "ApplyForCreateRoom"
)

type CreateRoomEvent struct {
	PlayerID  string // 玩家ID
	RoomTitle string // 房间名称
}
