package model

import (
	"fmt"
	"time"
)

// DummyCreatedAt is when the Dummy was created.
type DummyCreatedAt struct {
	value time.Time
}

// NewDummyCreatedAt rejects the zero time.
func NewDummyCreatedAt(t time.Time) (DummyCreatedAt, error) {
	if t.IsZero() {
		return DummyCreatedAt{}, fmt.Errorf("created at: zero time")
	}
	return DummyCreatedAt{value: t.UTC()}, nil
}

// Time returns the instant in UTC.
func (c DummyCreatedAt) Time() time.Time {
	return c.value
}
