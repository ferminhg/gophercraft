package event

import (
	"context"

	domainevent "github.com/fermin/gophercraft/internal/domain/event"
	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.EventPublisher = (*NoopPublisher)(nil)

// NoopPublisher implements port.EventPublisher by discarding all events (bootstrap / noop).
type NoopPublisher struct{}

// Publish implements port.EventPublisher.
func (NoopPublisher) Publish(context.Context, ...domainevent.DomainEvent) error {
	return nil
}
