package model

import (
	"fmt"
	"strings"
)

// DummyType classifies a Dummy variant.
type DummyType struct {
	value string
}

var (
	// DummyTypeAlpha is the alpha classification.
	DummyTypeAlpha = DummyType{value: "alpha"}
	// DummyTypeBeta is the beta classification.
	DummyTypeBeta = DummyType{value: "beta"}
	// DummyTypeGamma is the gamma classification.
	DummyTypeGamma = DummyType{value: "gamma"}
)

var validDummyTypeStrings = map[string]struct{}{
	DummyTypeAlpha.value: {},
	DummyTypeBeta.value:  {},
	DummyTypeGamma.value: {},
}

// NewDummyType parses and validates s as a known DummyType label.
func NewDummyType(s string) (DummyType, error) {
	trimmed := strings.TrimSpace(s)
	if _, ok := validDummyTypeStrings[trimmed]; !ok {
		return DummyType{}, fmt.Errorf("dummy type: invalid %q", trimmed)
	}
	return DummyType{value: trimmed}, nil
}

// String returns the canonical label.
func (t DummyType) String() string {
	return t.value
}

// IsValid reports whether t is a known DummyType value.
func (t DummyType) IsValid() bool {
	_, ok := validDummyTypeStrings[t.value]
	return ok
}
