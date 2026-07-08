package diff

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/martyn/apidiff/internal/model"
)

func compareEndpoints(oldSpec, newSpec *openapi3.T) []model.Change {
	var changes []model.Change

	oldEndpoints := collectEndpoints(oldSpec)
	newEndpoints := collectEndpoints(newSpec)

	// Detect removed endpoints
	for key := range oldEndpoints {
		if _, exists := newEndpoints[key]; !exists {
			method, path := splitEndpointKey(key)
			changes = append(changes, model.Change{
				Type:        model.EndpointRemoved,
				Severity:    model.SeverityBreaking,
				Category:    model.CategoryEndpoint,
				Path:        fmt.Sprintf("%s %s", method, path),
				Description: fmt.Sprintf("%s %s removed", method, path),
			})
		}
	}

	// Detect added endpoints
	for key := range newEndpoints {
		if _, exists := oldEndpoints[key]; !exists {
			method, path := splitEndpointKey(key)
			changes = append(changes, model.Change{
				Type:        model.EndpointAdded,
				Severity:    model.SeverityInfo,
				Category:    model.CategoryEndpoint,
				Path:        fmt.Sprintf("%s %s", method, path),
				Description: fmt.Sprintf("%s %s added", method, path),
			})
		}
	}

	return changes
}

// collectEndpoints builds a map of "METHOD:/path" -> true for all operations in the spec.
func collectEndpoints(spec *openapi3.T) map[string]bool {
	endpoints := make(map[string]bool)
	if spec.Paths == nil {
		return endpoints
	}

	for path, pathItem := range spec.Paths.Map() {
		for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
			op := getOperation(pathItem, method)
			if op != nil {
				key := method + ":" + path
				endpoints[key] = true
			}
		}
	}

	return endpoints
}

func getOperation(item *openapi3.PathItem, method string) *openapi3.Operation {
	switch method {
	case "GET":
		return item.Get
	case "POST":
		return item.Post
	case "PUT":
		return item.Put
	case "DELETE":
		return item.Delete
	case "PATCH":
		return item.Patch
	case "HEAD":
		return item.Head
	case "OPTIONS":
		return item.Options
	default:
		return nil
	}
}

func splitEndpointKey(key string) (method, path string) {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return key, ""
}
