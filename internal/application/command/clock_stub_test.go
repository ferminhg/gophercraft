package command_test

import "time"

// fixedClock always returns the same instant.
type fixedClock struct{ instant time.Time }

func (c fixedClock) Now() time.Time {
	return c.instant
}
