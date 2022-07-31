package msg

type CreateRoomRequest struct {
	RoomTitle string `json:"room_title"` // 房间名称
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"` // 房间ID
}

func NewCreateRoomRequest(roomTitle string) ([]byte, error) {
	// r := msgPack[CreateRoomRequest]{
	// 	Type: MSG_TYPE_CREATE_ROOM_REQUEST,
	// 	Payload: CreateRoomRequest{
	// 		RoomTitle: roomTitle,
	// 	},
	// }

	// bytes, err := json.Marshal(r)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
