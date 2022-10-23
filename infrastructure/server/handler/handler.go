package handler

import (
	"fmt"
	"net/http"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/gobang/manager"
	"github.com/gin-gonic/gin"
)

type HttpServerHandler struct {
	game *manager.GameManager
}

func NewHttpServerHandler(game *manager.GameManager) *HttpServerHandler {
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

func (handler *HttpServerHandler) CreateRoom() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 获取玩家信息
		id := ctx.Request.Header.Get("X-Game-Gobang-UserID")
		name := ctx.Request.Header.Get("X-Game-Gobang-UserName")

		fmt.Printf("[HTTP] 创建房间请求（id:%s, name:%s）", id, name)
		roomProfile, err := handler.game.CreateRoom(&entity.UserInfo{
			ID:         id,
			Username:   name,
			Nickname:   "",
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
