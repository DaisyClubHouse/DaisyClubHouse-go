package msg

import (
	"encoding/json"
	"log"
)

type msgPack[T interface{}] struct {
	Type    Kind `json:"type"`
	Payload T    `json:"payload"`
}

type Kind int

const (
	KindUnknown           Kind = iota // 未知消息
	KindCreateRoomRequest             // 创建房间
	KindCreateRoomResponse
)

func Parsing(data []byte) (Kind, []byte, error) {
	var packer msgPack[interface{}]
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
