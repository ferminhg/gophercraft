// Package port defines driven ports (repository and external service interfaces).
package port

import "time"

// Clock is a driven port for the current instant (testable replacement for time.Now).
type Clock interface {
	Now() time.Time
}
