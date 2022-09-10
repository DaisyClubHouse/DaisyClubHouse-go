package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerRoutes 注册HTTP路由
func registerRoutes(r *gin.Engine) {
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success!",
		})
	})
}
