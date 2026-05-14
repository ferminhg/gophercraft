package model

import "errors"

var (
	// ErrDummyIDEmpty indicates the DummyID string was blank after trimming.
	ErrDummyIDEmpty = errors.New("dummy id: empty")
	// ErrDummyIDInvalid indicates the DummyID string is not a valid UUID.
	ErrDummyIDInvalid = errors.New("dummy id: invalid uuid")
	// ErrDummyNameEmpty indicates the DummyName string was blank after trimming.
	ErrDummyNameEmpty = errors.New("dummy name: empty")
	// ErrDummyTypeInvalid indicates the DummyType string is not a known label.
	ErrDummyTypeInvalid = errors.New("dummy type: invalid")
	// ErrDummyCreatedAtZero indicates the creation time was zero.
	ErrDummyCreatedAtZero = errors.New("dummy created at: zero time")
)
