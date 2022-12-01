package room

import (
	"log"
	"sync"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/message/user_pack"
	"DaisyClubHouse/gobang/player"
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

	Owner       *player.Player
	Player      *player.Player
	lock        sync.Mutex
	whoseTurn   *player.Player
	whiteHolder *player.Player      // 执白棋的玩家（先行）
	blackHolder *player.Player      // 执黑棋的玩家（后行）
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

func (room *Room) PlayerJoin(playerInfo *entity.UserInfo, client *player.Player) {
	room.lock.Lock()
	defer room.lock.Unlock()

	client.UserInfo = playerInfo

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

func (room *Room) PlayerPlaceThePiece(playerID string, x, y int) {
	if room.Status != Status_Playing {
		slog.Info("游戏未开始", slog.String("player_id", playerID))
		return
	}

	if room.whoseTurn != nil && room.whoseTurn.ID != playerID {
		slog.Info("不是你的回合", slog.String("player_id", playerID))
		return
	}

	if room.whiteHolder.ID == playerID {
		room.whiteMatrix.Put(x, y)
	} else if room.blackHolder.ID == playerID {
		room.blackMatrix.Put(x, y)
	} else {
		slog.Info("不是你的回合", slog.String("player_id", playerID))
		return
	}

	// 广播玩家落子信息
	pack := user_pack.UserPack[user_pack.BroadcastPlayerPlaceThePiece]{
		Type: user_pack.KindBroadcastPlayerPlaceThePiece,
		Payload: user_pack.BroadcastPlayerPlaceThePiece{
			RoomID:     room.ID,
			PlayerID:   room.whoseTurn.PlayerID(),
			PieceWhite: room.whiteHolder.ID == room.whoseTurn.ID,
			X:          x,
			Y:          y,
		},
	}

	go room.Broadcast(pack.Marshal())

	go room.turnChanged(true)

	// 判断是否获胜
	if room.whiteMatrix.IsWin() {
		room.gameSettlement(room.whiteHolder.PlayerProfile())
		slog.Info("白棋获胜", slog.Any("winner", room.whiteHolder.PlayerProfile()))
		return
	} else if room.blackMatrix.IsWin() {
		room.gameSettlement(room.whiteHolder.PlayerProfile())
		slog.Info("黑棋获胜", slog.Any("winner", room.blackHolder.PlayerProfile()))
		return
	}
}

func (room *Room) gameSettlement(winner *user_pack.PlayerProfile) {
	// 结算中
	room.Status = Status_Settlement

	// 广播游戏结算消息
	pack := user_pack.UserPack[user_pack.BroadcastGameOver]{
		Type: user_pack.KindBroadcastGameOver,
		Payload: user_pack.BroadcastGameOver{
			Winner: winner,
		},
	}

	go room.Broadcast(pack.Marshal())
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

	pack := user_pack.UserPack[user_pack.BroadcastGameBeginning]{
		Type: user_pack.KindBroadcastRoomGameBeginning,
		Payload: user_pack.BroadcastGameBeginning{
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

	pack := user_pack.UserPack[user_pack.BroadcastPlayerAction]{
		Type: user_pack.KindBroadcastPlayerAction,
		Payload: user_pack.BroadcastPlayerAction{
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
