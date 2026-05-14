package command

import "errors"

// ErrInvalidDummy is returned when Dummy fails domain validation.
var ErrInvalidDummy = errors.New("invalid dummy")
