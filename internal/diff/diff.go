package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/martyn/apidiff/internal/model"
)

// Compare takes two OpenAPI 3.x specs and returns a list of changes between them.
func Compare(oldSpec, newSpec *openapi3.T) []model.Change {
	var changes []model.Change

	changes = append(changes, compareEndpoints(oldSpec, newSpec)...)
	changes = append(changes, compareSchemas(oldSpec, newSpec)...)

	return changes
}

// Summary holds counts of changes by type.
type Summary struct {
	Added         int
	Removed       int
	SchemaChanges int
	Breaking      int
}

// Summarize computes summary statistics from a list of changes.
func Summarize(changes []model.Change) Summary {
	var s Summary
	for _, c := range changes {
		switch c.Type {
		case model.EndpointAdded:
			s.Added++
		case model.EndpointRemoved:
			s.Removed++
		default:
			s.SchemaChanges++
		}
		if c.IsBreaking() {
			s.Breaking++
		}
	}
	return s
}
