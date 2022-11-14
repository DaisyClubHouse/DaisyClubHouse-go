package room

import (
	"log"
	"sync"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/manager/player"
	"DaisyClubHouse/gobang/msg"
	"DaisyClubHouse/utils"
	"golang.org/x/exp/slog"
)

// RoomProfile 房间概要信息
type RoomProfile struct {
	ID         string           `json:"id"`         // 房间唯一ID
	Title      string           `json:"title"`      // 房间名称
	Status     Status           `json:"status"`     // 房间状态
	Creator    *entity.UserInfo `json:"creator"`    // 房主信息
	Player     *entity.UserInfo `json:"player"`     // 玩家信息
	CreateTime time.Time        `json:"createTime"` // 创建时间
}

// Room 房间
type Room struct {
	RoomProfile // 房间概要信息

	Owner       *player.Client
	Player      *player.Client
	lock        sync.Mutex
	whoseTurn   *player.Client
	whiteHolder *player.Client      // 执白棋的玩家（先行）
	blackHolder *player.Client      // 执黑棋的玩家（后行）
	whiteMatrix *entity.ChessMatrix // 白棋的棋盘
	blackMatrix *entity.ChessMatrix // 黑棋的棋盘
}

type Status int

const (
	Status_Waiting    Status = iota // 等待玩家加入
	Status_Playing                  // 游戏中
	Status_Settlement               // 游戏结算中
	Status_Zombie                   // 游戏结束待回收
)

// CreateNewRoom 创建新房间
func CreateNewRoom(owner *entity.UserInfo) *Room {
	profile := RoomProfile{
		ID:         utils.GenerateRandomID(),
		Title:      owner.Username + "的房间",
		Status:     Status_Waiting,
		Creator:    owner,
		Player:     nil,
		CreateTime: time.Now(),
	}
	room := &Room{
		RoomProfile: profile,
		Owner:       nil,
		Player:      nil,
		whoseTurn:   nil,
		whiteHolder: nil,
		blackHolder: nil,
		whiteMatrix: entity.NewChessMatrix(15),
		blackMatrix: entity.NewChessMatrix(15),
	}

	return room
}

func (room *Room) PlayerJoin(playerInfo *entity.UserInfo, client *player.Client) {
	room.lock.Lock()
	defer room.lock.Unlock()

	client.Identity = playerInfo

	if room.RoomProfile.Creator.ID == playerInfo.ID {
		// 房主
		room.Owner = client
		slog.Info("房主客户端连接上房间",
			slog.String("room_id", room.ID),
			slog.String("client_id", client.ID),
			slog.Any("player_id", playerInfo),
		)
	} else {
		room.Player = client
		slog.Info("玩家客户端连接上房间",
			slog.String("room_id", room.ID),
			slog.String("client_id", client.ID),
			slog.Any("player_id", playerInfo),
		)
	}

	// 状态流转
	if room.Owner != nil && room.Player != nil {
		room.gameBegin()
	}
}

func (room *Room) OwnerHold() string {
	if room.whiteHolder == room.Owner {
		return "执白"
	} else if room.blackHolder == room.Owner {
		return "执黑"
	} else {
		return "未知"
	}
}

func (room *Room) PlayerHold() string {
	if room.whiteHolder == room.Player {
		return "执白"
	} else if room.blackHolder == room.Player {
		return "执黑"
	} else {
		return "未知"
	}
}

func (room *Room) PlayerPlaceThePiece(playerID string, x, y int) string {
	if room.Status != Status_Playing {
		return "游戏未开始"
	}

	if room.whoseTurn != nil && room.whoseTurn.ID != playerID {
		return "不是你的回合"
	}

	if room.whiteHolder.ID == playerID {
		room.whiteMatrix.Put(x, y)
	} else if room.blackHolder.ID == playerID {
		room.blackMatrix.Put(x, y)
	} else {
		return "不是你的回合"
	}

	// 判断是否获胜
	// if room.whiteMatrix.IsWin(x, y) {
	// 	room.status = Status_Settlement
	// 	return "白棋获胜"
	// } else if room.blackMatrix.IsWin(x, y) {
	// 	room.status = Status_Settlement
	// 	return "黑棋获胜"
	// }

	// 广播玩家落子信息
	pack := msg.UserPack[msg.BroadcastPlayerPlaceThePiece]{
		Type: msg.KindBroadcastPlayerPlaceThePiece,
		Payload: msg.BroadcastPlayerPlaceThePiece{
			RoomID:     room.ID,
			PlayerID:   room.whoseTurn.PlayerID(),
			PieceWhite: room.whiteHolder.ID == room.whoseTurn.ID,
			X:          x,
			Y:          y,
		},
	}

	go room.Broadcast(pack.Marshal())

	go room.turnChanged(true)
	return "OK"
}

func (room *Room) gameBegin() {
	room.Status = Status_Playing

	// 准备开始游戏
	log.Println("----------Ready to start game----------")

	// 随机决定执黑执白
	if utils.RandomHalfRate() {
		room.whiteHolder = room.Player
		room.blackHolder = room.Owner
	} else {
		room.whiteHolder = room.Owner
		room.blackHolder = room.Player
	}

	log.Printf("|Room[%s]", room.ID)
	log.Printf("|Owner%s[%s,%s]",
		room.OwnerHold(),
		room.Owner.ID, room.Owner.RemoteAddr())
	log.Printf("|Player%s[%s,%s]",
		room.PlayerHold(),
		room.Player.ID, room.Player.RemoteAddr())

	// 初始化通知执白执黑已经谁先行
	room.whoseTurn = room.blackHolder

	pack := msg.UserPack[msg.BroadcastGameBeginning]{
		Type: msg.KindBroadcastRoomGameBeginning,
		Payload: msg.BroadcastGameBeginning{
			RoomID:      room.ID,
			WhiteHolder: room.whiteHolder.PlayerProfile(),
			BlackHolder: room.blackHolder.PlayerProfile(),
		},
	}
	room.Broadcast(pack.Marshal())

	room.turnChanged(false)
}

// DisconnectGameOver 玩家断线，游戏结束
func (room *Room) DisconnectGameOver() {
	slog.Info("玩家断线游戏结束", slog.String("room_id", room.ID))

	if room.RoomProfile.Status == Status_Zombie {
		return
	}

	// 发送公告
	if p := room.Owner; p != nil {
		p.NoticedDisconnect()
	}
	if p := room.Player; p != nil {
		p.NoticedDisconnect()
	}

	room.RoomProfile.Status = Status_Zombie
}

func (room *Room) turnChanged(turn bool) {
	slog.Info("回合变更", slog.Bool("turn", turn))

	if turn {
		if room.whoseTurn == room.whiteHolder {
			room.whoseTurn = room.blackHolder
		} else {
			room.whoseTurn = room.whiteHolder
		}
	}

	pack := msg.UserPack[msg.BroadcastPlayerAction]{
		Type: msg.KindBroadcastPlayerAction,
		Payload: msg.BroadcastPlayerAction{
			RoomId:    room.ID,
			WhoseTurn: room.whoseTurn.PlayerProfile(),
		},
	}
	room.Broadcast(pack.Marshal())

	slog.Info("广播玩家行动", slog.Bool("turn", turn))
}

// Broadcast 房间内广播
func (room *Room) Broadcast(data []byte) {
	log.Printf("[广播] room:%s lens:%d", room.ID, len(data))
	slog.Info("广播房间消息", slog.Any("data", string(data)))
	if len(data) == 0 {
		// 如果没有数据，则不发送
		log.Println("------------------ERROR------------------")
		log.Printf("[ERROR] Room[%s] Broadcast data is empty\n", room.ID)
		return
	}

	if room.Player != nil {
		room.Player.SendRawMessage(data)
	}
	if room.Owner != nil {
		room.Owner.SendRawMessage(data)
	}
}
