package events

import (
	"sync"
)

// EventBus publishes events to subscribers without direct coupling between components.
// See ADR-003 and tasks/004-event-bus.md.
type EventBus interface {
	Publish(event Event)
	Subscribe(eventType Type, handler func(Event)) (unsubscribe func())
}

// ChannelBus is an in-process EventBus backed by buffered channels and goroutine dispatch.
type ChannelBus struct {
	mu          sync.RWMutex
	subscribers map[Type][]chan Event
	bufferSize  int
}

// NewChannelBus returns a channel-based event bus with the given per-subscriber buffer size.
func NewChannelBus(bufferSize int) *ChannelBus {
	if bufferSize < 1 {
		bufferSize = 16
	}
	return &ChannelBus{
		subscribers: make(map[Type][]chan Event),
		bufferSize:  bufferSize,
	}
}

// Publish delivers the event to all subscribers of its type without blocking the caller.
func (b *ChannelBus) Publish(event Event) {
	b.mu.RLock()
	subs := append([]chan Event(nil), b.subscribers[event.Type]...)
	b.mu.RUnlock()

	for _, ch := range subs {
		ch := ch
		go func() {
			ch <- event
		}()
	}
}

// Subscribe registers a handler for eventType. Returns an unsubscribe function.
func (b *ChannelBus) Subscribe(eventType Type, handler func(Event)) (unsubscribe func()) {
	ch := make(chan Event, b.bufferSize)

	b.mu.Lock()
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	b.mu.Unlock()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					return
				}
				handler(ev)
			case <-done:
				return
			}
		}
	}()

	return func() {
		close(done)
		b.mu.Lock()
		defer b.mu.Unlock()

		subs := b.subscribers[eventType]
		for i, sub := range subs {
			if sub == ch {
				close(ch)
				b.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}
