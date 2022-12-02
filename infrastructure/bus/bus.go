package bus

import "github.com/asaskevich/EventBus"

func EventBusProvider() EventBus.Bus {
	bus := EventBus.New()

	return bus
}
