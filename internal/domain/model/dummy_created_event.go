package model

import (
	"time"

	"github.com/fermin/gophercraft/internal/domain/event"
)

// DummyCreated is emitted when a new Dummy aggregate is created.
type DummyCreated struct {
	dummyID      string
	name         string
	dummyType    string
	dummyCreated time.Time
	occurredAt   time.Time
}

// NewDummyCreated builds a DummyCreated event from validated value objects.
func NewDummyCreated(id DummyID, name DummyName, t DummyType, createdAt DummyCreatedAt) *DummyCreated {
	tm := createdAt.Time()
	return &DummyCreated{
		dummyID:      id.String(),
		name:         name.String(),
		dummyType:    t.String(),
		dummyCreated: tm,
		occurredAt:   tm,
	}
}

var _ event.DomainEvent = (*DummyCreated)(nil)

// EventName returns the stable event identifier.
func (e *DummyCreated) EventName() string {
	return "dummy.created"
}

// AggregateID returns the Dummy aggregate UUID string.
func (e *DummyCreated) AggregateID() string {
	return e.dummyID
}

// OccurredAt is when the event was recorded (same instant as creation by default).
func (e *DummyCreated) OccurredAt() time.Time {
	return e.occurredAt
}

// DummyID is the persisted aggregate identifier.
func (e *DummyCreated) DummyID() string {
	return e.dummyID
}

// DummyName is the persisted name snapshot.
func (e *DummyCreated) DummyName() string {
	return e.name
}

// DummyType is the persisted type label.
func (e *DummyCreated) DummyType() string {
	return e.dummyType
}

// DummyCreatedAt is the creation timestamp captured in the event.
func (e *DummyCreated) DummyCreatedAt() time.Time {
	return e.dummyCreated
}
