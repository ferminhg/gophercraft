package port

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/event"
)

// EventPublisher publishes domain events after successful persistence (e.g. bus, log, outbox).
type EventPublisher interface {
	Publish(ctx context.Context, events ...event.DomainEvent) error
}
