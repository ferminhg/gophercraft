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
	logger  port.Logger
	metrics port.MetricsRecorder
}

// NewNoopPublisher returns a publisher that logs and discards events.
func NewNoopPublisher(logger port.Logger, metrics port.MetricsRecorder) *NoopPublisher {
	return &NoopPublisher{logger: logger, metrics: metrics}
}

// Publish implements port.EventPublisher.
func (p *NoopPublisher) Publish(_ context.Context, events ...domainevent.DomainEvent) error {
	for i := range events {
		eventName := events[i].EventName()
		p.metrics.RecordEventPublished(eventName)
		p.logger.Info("noop publisher received event", "event", eventName, "detail", events[i].String())
	}
	return nil
}
