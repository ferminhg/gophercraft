package model

import "strings"

// DummyNameMaxLength is the maximum allowed length for a DummyName after trimming.
const DummyNameMaxLength = 255

// DummyName is the display name of a Dummy.
type DummyName struct {
	value string
}

// NewDummyName returns a trimmed non-empty DummyName.
func NewDummyName(s string) (*DummyName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, ErrDummyNameEmpty
	}
	if len(s) > DummyNameMaxLength {
		return nil, ErrDummyNameTooLong
	}
	return &DummyName{value: s}, nil
}

// String returns the name value.
func (n DummyName) String() string {
	return n.value
}
