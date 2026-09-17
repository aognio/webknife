package webknife

import "github.com/webknife/webknife/internal/events"

type EventPublisher struct {
	bus *events.EventBus
}

func NewEventPublisher(bus *events.EventBus) *EventPublisher {
	return &EventPublisher{bus: bus}
}

func (p *EventPublisher) Publish(event events.Event) {
	p.bus.Publish(event)
}
