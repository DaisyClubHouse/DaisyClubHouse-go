package hub

import (
	"context"

	"go.uber.org/fx"
)

func GameHubProvider(lc fx.Lifecycle) *GameHub {
	hub := NewGameHub()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go hub.run()

			return nil
		},
		OnStop: nil,
	})

	return hub
}
