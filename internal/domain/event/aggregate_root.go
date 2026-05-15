package event

// AggregateRoot holds unpublished domain events for an aggregate root.
// Embed zero-value AggregateRoot into entities that raise events.
type AggregateRoot struct {
	events []DomainEvent
}

// RecordEvent appends an event to the unpublished buffer.
func (r *AggregateRoot) RecordEvent(e DomainEvent) {
	r.events = append(r.events, e)
}

// PullDomainEvents returns and clears unpublished events.
func (r *AggregateRoot) PullDomainEvents() []DomainEvent {
	out := r.events
	r.events = nil
	return out
}

// ClearDomainEvents discards unpublished events without returning them.
func (r *AggregateRoot) ClearDomainEvents() {
	r.events = nil
}
