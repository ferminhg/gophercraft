// Package port defines driven ports (repository and external service interfaces).
package port

// UUIDGenerator is a driven port for generating unique identifiers.
type UUIDGenerator interface {
	Generate() string
}
