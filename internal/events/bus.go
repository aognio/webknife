package events

import (
	"sync"
)

type EventBus struct {
	subscribers []func(Event)
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subscribers {
		sub(event)
	}
}

func (b *EventBus) Subscribe(fn func(Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, fn)
}
