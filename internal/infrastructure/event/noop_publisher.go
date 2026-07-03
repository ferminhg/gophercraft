// Package event provides infrastructure adapters for publishing domain events.
package event

import (
	"context"

	domainevent "github.com/fermin/gophercraft/internal/domain/event"
	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.EventPublisher = (*NoopPublisher)(nil)

// NoopPublisher implements port.EventPublisher by discarding all events (bootstrap / noop).
type NoopPublisher struct {
	logger port.Logger
}

// NewNoopPublisher returns a publisher that logs and discards events.
func NewNoopPublisher(logger port.Logger) *NoopPublisher {
	return &NoopPublisher{logger: logger}
}

// Publish implements port.EventPublisher.
func (p *NoopPublisher) Publish(_ context.Context, events ...domainevent.DomainEvent) error {
	for i := range events {
		p.logger.Info("noop publisher received event", "event", events[i].EventName(), "detail", events[i].String())
	}
	return nil
}
