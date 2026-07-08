package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/martyn/apidiff/internal/model"
)

func compareSchemas(oldSpec, newSpec *openapi3.T) []model.Change {
	var changes []model.Change

	oldSchemas := collectSchemas(oldSpec)
	newSchemas := collectSchemas(newSpec)

	for name, oldSchema := range oldSchemas {
		newSchema, exists := newSchemas[name]
		if !exists {
			continue
		}
		changes = append(changes, compareSchemaProperties(name, oldSchema, newSchema)...)
	}

	return changes
}

func compareSchemaProperties(schemaName string, oldSchema, newSchema *openapi3.SchemaRef) []model.Change {
	var changes []model.Change

	if oldSchema.Value == nil || newSchema.Value == nil {
		return changes
	}

	oldProps := oldSchema.Value.Properties
	newProps := newSchema.Value.Properties

	// Detect removed properties
	for propName := range oldProps {
		if _, exists := newProps[propName]; !exists {
			changes = append(changes, model.Change{
				Type:        model.PropertyRemoved,
				Severity:    model.SeverityBreaking,
				Category:    model.CategoryProperty,
				Path:        schemaName,
				Property:    fmt.Sprintf("%s.%s", schemaName, propName),
				Description: fmt.Sprintf("Property %s removed from %s", propName, schemaName),
			})
		}
	}

	// Detect added properties
	for propName := range newProps {
		if _, exists := oldProps[propName]; !exists {
			severity := model.SeverityInfo
			// Adding a required property is breaking
			if isRequired(newSchema.Value, propName) {
				severity = model.SeverityBreaking
			}
			changes = append(changes, model.Change{
				Type:        model.PropertyAdded,
				Severity:    severity,
				Category:    model.CategoryProperty,
				Path:        schemaName,
				Property:    fmt.Sprintf("%s.%s", schemaName, propName),
				Description: fmt.Sprintf("Property %s added to %s", propName, schemaName),
			})
		}
	}

	// Detect type changes on existing properties
	for propName, oldProp := range oldProps {
		newProp, exists := newProps[propName]
		if !exists {
			continue
		}
		if oldProp.Value != nil && newProp.Value != nil {
			if oldProp.Value.Type.Slice() != nil && newProp.Value.Type.Slice() != nil {
				oldType := typeString(oldProp.Value)
				newType := typeString(newProp.Value)
				if oldType != newType {
					changes = append(changes, model.Change{
						Type:        model.PropertyTypeChanged,
						Severity:    model.SeverityBreaking,
						Category:    model.CategoryProperty,
						Path:        schemaName,
						Property:    fmt.Sprintf("%s.%s", schemaName, propName),
						Description: fmt.Sprintf("Property %s type changed from %s to %s in %s", propName, oldType, newType, schemaName),
					})
				}
			}
		}
	}

	return changes
}

func collectSchemas(spec *openapi3.T) map[string]*openapi3.SchemaRef {
	schemas := make(map[string]*openapi3.SchemaRef)
	if spec.Components == nil || spec.Components.Schemas == nil {
		return schemas
	}
	for name, schema := range spec.Components.Schemas {
		schemas[name] = schema
	}
	return schemas
}

func isRequired(schema *openapi3.Schema, propName string) bool {
	for _, r := range schema.Required {
		if r == propName {
			return true
		}
	}
	return false
}

func typeString(schema *openapi3.Schema) string {
	types := schema.Type.Slice()
	if len(types) == 0 {
		return "unknown"
	}
	return types[0]
}
