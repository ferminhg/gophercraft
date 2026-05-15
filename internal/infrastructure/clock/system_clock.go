package clock

import (
	"time"

	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.Clock = (*SystemClock)(nil)

// SystemClock returns time.Now in UTC on each Now() call.
type SystemClock struct{}

// Now implements port.Clock.
func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
