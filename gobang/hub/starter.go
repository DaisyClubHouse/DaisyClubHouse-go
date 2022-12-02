package hub

import (
	"context"

	"go.uber.org/fx"
)

func GameHubInvoker(lc fx.Lifecycle, hub *GameHub) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go hub.run()

			return nil
		},
		OnStop: nil,
	})

	return nil
}
