package server

import (
	"log"
	"net/http"
	"time"

	"DaisyClubHouse/domain/entity"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// registerRoutes 注册HTTP路由
func registerRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success!",
		})
	})
	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})

	// api
	// api := r.Group("/api", authMiddleware.MiddlewareFunc())
	api := r.Group("/api")
	{
		api.GET("/refresh_token", authMiddleware.RefreshHandler)
		api.GET("/hello", helloHandler)

		api.GET("/game/room/list", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": gin.H{
					"list": []entity.RoomProfile{
						{
							ID:         "1",
							Title:      "房间0001",
							Status:     0,
							CreateTime: time.Now(),
						},
						{
							ID:         "2",
							Title:      "房间0002",
							Status:     0,
							CreateTime: time.Now(),
						},
					},
				},
			})
		})
	}
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Println(claims)
	c.JSON(200, gin.H{
		"hello": "world",
	})
}
