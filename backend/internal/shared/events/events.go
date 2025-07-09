package events

import (
	"context"
	"reflect"
	"sync"
	"time"
)

// EventType represents the type of an event
type EventType string

// Event represents a domain event
type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Source    string      `json:"source"`
	Timestamp time.Time   `json:"timestamp"`
	Metadata  Metadata    `json:"metadata"`
}

// Metadata contains additional information about an event
type Metadata struct {
	UserID             string            `json:"user_id,omitempty"`
	ProviderIdentifier string            `json:"provider_identifier,omitempty"`
	Operation          string            `json:"operation,omitempty"`
	TraceID            string            `json:"trace_id,omitempty"`
	SpanID             string            `json:"span_id,omitempty"`
	Properties         map[string]string `json:"properties,omitempty"`
}

// NewEvent creates a new event
func NewEvent(eventType EventType, payload interface{}, metadata Metadata) Event {
	return Event{
		Type:      string(eventType),
		Payload:   payload,
		Source:    "context-space-backend",
		Timestamp: time.Now(),
		Metadata:  metadata,
	}
}

// EventHandler is a function that handles events
type EventHandler func(ctx context.Context, event Event) error

// Bus is an event bus for publishing and subscribing to events
type Bus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex // Protects subscribers map
}

// NewBus creates a new event bus
func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (b *Bus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

// Publish publishes an event to all subscribers
func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, ok := b.subscribers[event.Type]
	if !ok {
		b.mu.RUnlock()
		// No subscribers for this event type
		return nil
	}

	// Make a copy to avoid holding lock during handler execution
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	for _, handler := range handlersCopy {
		if err := handler(ctx, event); err != nil {
			// Log the error but continue with other handlers
			return err
		}
	}

	return nil
}

// PublishDomainEvent implements the EventPublisher interface for domain events
func (b *Bus) PublishDomainEvent(event interface{}) {
	// Create a background context for event publishing
	ctx := context.Background()

	// For domain events, we need to wrap them in our Event type
	// First, determine the event type from the concrete type name
	eventType := reflect.TypeOf(event).Elem().Name()

	// Create our event wrapper
	e := Event{
		Type:      eventType,
		Payload:   event,
		Source:    "context-space-backend",
		Timestamp: time.Now(),
		Metadata:  Metadata{},
	}

	// Publish the event
	b.Publish(ctx, e)
}
