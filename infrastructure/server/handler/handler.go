package handler

import (
	"fmt"
	"net/http"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/manager/game"
	"github.com/gin-gonic/gin"
)

type HttpServerHandler struct {
	game *game.GameManager
}

func NewHttpServerHandler(game *game.GameManager) *HttpServerHandler {
	return &HttpServerHandler{game: game}
}

func (handler *HttpServerHandler) GetRoomProfileListHandler() func(*gin.Context) {
	return func(ctx *gin.Context) {
		roomProfileList := handler.game.RoomProfileList()

		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"list": roomProfileList,
			},
		})
	}
}

type CreateRoomRequest struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

func (handler *HttpServerHandler) CreateRoom() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 获取玩家信息
		var req CreateRoomRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "请求参数绑定错误",
			})
			return
		}
		// 提取玩家ID，name
		id, name := req.UserID, req.UserName

		fmt.Printf("[HTTP] 创建房间请求（player_id:%s, player_name:%s）", id, name)
		roomProfile, err := handler.game.CreateRoom(&entity.UserInfo{
			ID:         id,
			Username:   name,
			Nickname:   "",
			Avatar:     fmt.Sprintf("https://joeschmoe.io/api/v1/%s", name),
			CreateTime: time.Now(),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"room_id":    roomProfile.ID,
				"room_title": roomProfile.Title,
			},
		})
	}
}
