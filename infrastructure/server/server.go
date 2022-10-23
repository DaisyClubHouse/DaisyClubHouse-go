package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"DaisyClubHouse/infrastructure/server/handler"
	"DaisyClubHouse/infrastructure/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const httpAddr = ":3000"

// HttpServer 提供HTTP服务
func HttpServer(lc fx.Lifecycle, wsFunc gin.HandlerFunc, handler *handler.HttpServerHandler) {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	// ws
	r.GET("/ws", wsFunc)

	// http
	jwt := middleware.Authorization()
	RegisterRoutes(r, jwt, handler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := r.Run(httpAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
					log.Printf("listen: %s\n", err)
				}
			}()

			return nil
		},
		OnStop: nil,
	})
}
