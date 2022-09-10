package server

import (
	"log"
	"net/http"

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

	// auth
	auth := r.Group("/auth", authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	// api
	api := r.Group("/api", authMiddleware.MiddlewareFunc())
	{
		api.GET("/hello", helloHandler)
	}
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Println(claims)
	c.JSON(200, gin.H{
		"hello": "world",
	})
}
