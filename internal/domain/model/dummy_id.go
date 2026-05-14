package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// DummyID identifies a Dummy aggregate instance.
type DummyID struct {
	value string
}

// NewDummyID parses and validates s as a UUID.
func NewDummyID(s string) (*DummyID, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil, ErrDummyIDEmpty
	}
	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDummyIDInvalid, err)
	}
	return &DummyID{value: parsed.String()}, nil
}

// String returns the canonical UUID string form.
func (id DummyID) String() string {
	return id.value
}
