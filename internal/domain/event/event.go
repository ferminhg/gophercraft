// Package event defines domain-event contracts shared by aggregates.
package event

import "time"

// DomainEvent is a domain-level fact that occurred in the past.
type DomainEvent interface {
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}
