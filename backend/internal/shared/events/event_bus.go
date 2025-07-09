package events

import (
	"context"
)

// EventBus defines the interface for event publishing and subscription
type EventBus interface {
	// Publish publishes an event to all subscribers
	Publish(ctx context.Context, event Event) error

	// Subscribe registers a handler for a specific event type
	Subscribe(eventType string, handler EventHandler)

	// Unsubscribe removes a handler for a specific event type
	Unsubscribe(eventType string, handler EventHandler)
}

// EventPublisher is responsible for publishing events
type EventPublisher interface {
	// PublishDomainEvent publishes a domain event
	PublishDomainEvent(event interface{})
}

// Ensure Bus implements EventBus
var _ EventBus = (*Bus)(nil)

// Ensure Bus implements EventPublisher
var _ EventPublisher = (*Bus)(nil)

// Unsubscribe removes a handler for a specific event type
func (b *Bus) Unsubscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	handlers, ok := b.subscribers[eventType]
	if !ok {
		return
	}

	for i, h := range handlers {
		if &h == &handler {
			// Remove the handler by slicing
			b.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	// If no more handlers, remove the event type
	if len(b.subscribers[eventType]) == 0 {
		delete(b.subscribers, eventType)
	}
}

// Payload is a generic map for event data
type Payload map[string]interface{}
