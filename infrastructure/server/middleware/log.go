package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

// Logger HTTP日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		slog.Info("Request received",
			slog.String("client_ip", c.ClientIP()),
			slog.String("proto", c.Request.Proto),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)
		c.Next()

		slog.Info("Response completed",
			slog.String("client_ip", c.ClientIP()),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(startTime)),
		)
	}
}
