package model

// ChangeType represents the kind of API change detected.
type ChangeType string

const (
	EndpointAdded       ChangeType = "endpoint_added"
	EndpointRemoved     ChangeType = "endpoint_removed"
	SchemaChanged       ChangeType = "schema_changed"
	PropertyAdded       ChangeType = "property_added"
	PropertyRemoved     ChangeType = "property_removed"
	PropertyTypeChanged ChangeType = "property_type_changed"
)

// Severity indicates how impactful a change is.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityBreaking Severity = "breaking"
)

// Category groups changes by area of the spec.
type Category string

const (
	CategoryEndpoint  Category = "endpoint"
	CategorySchema    Category = "schema"
	CategoryProperty  Category = "property"
	CategoryParameter Category = "parameter"
	CategoryResponse  Category = "response"
)

// Change represents a single difference between two API specifications.
type Change struct {
	Type        ChangeType `json:"type"`
	Severity    Severity   `json:"severity"`
	Category    Category   `json:"category"`
	Path        string     `json:"path"`
	Property    string     `json:"property,omitempty"`
	Description string     `json:"description"`
}

// IsBreaking returns true if this change is a breaking change.
func (c Change) IsBreaking() bool {
	return c.Severity == SeverityBreaking
}
