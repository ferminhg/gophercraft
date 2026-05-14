package model

import (
	"fmt"
	"strings"
)

// DummyName is the display name of a Dummy.
type DummyName struct {
	value string
}

// NewDummyName returns a trimmed non-empty DummyName.
func NewDummyName(s string) (DummyName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return DummyName{}, fmt.Errorf("dummy name: empty")
	}
	return DummyName{value: s}, nil
}

// String returns the name value.
func (n DummyName) String() string {
	return n.value
}
