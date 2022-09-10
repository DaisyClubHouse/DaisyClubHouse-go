package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const httpAddr = ":3000"

// ServerProvider HTTP服务
func ServerProvider() {
	r := gin.Default()
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success!",
		})
	})
	r.Run(httpAddr)
}
