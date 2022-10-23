package handler

import (
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
