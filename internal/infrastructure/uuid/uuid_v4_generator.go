// Package uuid provides infrastructure adapters for UUID generation.
package uuid

import "github.com/google/uuid"

// GoogleUUIDGenerator implements port.UUIDGenerator using github.com/google/uuid.
type GoogleUUIDGenerator struct{}

// Generate returns a new random UUID v4 string.
func (GoogleUUIDGenerator) Generate() string {
	return uuid.NewString()
}
