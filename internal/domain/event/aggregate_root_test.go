package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	domainevent "github.com/fermin/gophercraft/internal/domain/event"
)

type stubDomainEvent struct {
	name        string
	aggregateID string
	occurredAt  time.Time
}

func (e stubDomainEvent) EventName() string     { return e.name }
func (e stubDomainEvent) AggregateID() string   { return e.aggregateID }
func (e stubDomainEvent) OccurredAt() time.Time { return e.occurredAt }

func TestAggregateRoot_RecordAndPull_ReturnsEventsAndClears(t *testing.T) {
	t.Parallel()

	var root domainevent.AggregateRoot
	now := time.Date(2026, 5, 14, 10, 0, 0, 0, time.UTC)
	ev1 := stubDomainEvent{name: "a", aggregateID: "id1", occurredAt: now}
	ev2 := stubDomainEvent{name: "b", aggregateID: "id1", occurredAt: now.Add(time.Minute)}
	root.RecordEvent(ev1)
	root.RecordEvent(ev2)

	pulled := root.PullDomainEvents()
	require.Len(t, pulled, 2)
	require.Equal(t, "a", pulled[0].EventName())
	require.Equal(t, "b", pulled[1].EventName())

	pulledAgain := root.PullDomainEvents()
	require.Empty(t, pulledAgain)
}

func TestAggregateRoot_ClearDomainEvents_EmptiesWithoutReturn(t *testing.T) {
	t.Parallel()

	var root domainevent.AggregateRoot
	root.RecordEvent(stubDomainEvent{name: "only", aggregateID: "x", occurredAt: time.Now().UTC()})
	root.ClearDomainEvents()

	require.Empty(t, root.PullDomainEvents())
}

func TestAggregateRoot_ZeroRecordAndPull_NoPanic(t *testing.T) {
	t.Parallel()

	var root domainevent.AggregateRoot
	require.Empty(t, root.PullDomainEvents())
}
