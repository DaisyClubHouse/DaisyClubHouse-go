package msg

import (
	"encoding/json"
	"log"
)

type UserPack[T interface{}] struct {
	Type    Kind `json:"type"`
	Payload T    `json:"payload"`
}

func (p *UserPack[T]) Marshal() []byte {
	bytes, _ := json.Marshal(p)
	return bytes
}

type Kind int

const (
	KindUnknown               Kind = iota // 未知消息
	KindJoinRoomRequest            = 10   // 请求加入房间
	KindJoinRoomResponse           = 11   //
	KindPlaceThePieceRequest       = 20   // 请求下棋
	KindPlaceThePieceResponse      = 21   //
)
const (
	KindBroadcastRoomGameBeginning   Kind = iota + 100 // 广播游戏开始 100
	KindBroadcastPlayerPlaceThePiece                   // 广播玩家操作落子
)

func Parsing(data []byte) (Kind, []byte, error) {
	var packer UserPack[interface{}]
	if err := json.Unmarshal(data, &packer); err != nil {
		return KindUnknown, nil, err
	}

	log.Println("packer: {$v}", packer)

	bytes, err := json.Marshal(packer.Payload)
	if err != nil {
		return KindUnknown, nil, err
	}

	return packer.Type, bytes, nil
}
