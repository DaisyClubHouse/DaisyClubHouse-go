// Code generated by hertz generator.

package main

import (
	"DaisyClubHouse/infrastructure/http"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Invoke(http.HertzProvider),
	)
	app.Run()
}
