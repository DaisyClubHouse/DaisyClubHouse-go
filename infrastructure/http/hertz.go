package http

import "github.com/cloudwego/hertz/pkg/app/server"

func HertzProvider() {
	h := server.Default()

	// middleware
	h.Use(accessLog())

	// register routers
	register(h)

	// run
	h.Spin()
}
