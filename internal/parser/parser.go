package parser

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
)

// Parse loads an OpenAPI specification from a file path.
// It supports both OpenAPI 2.x (Swagger) and 3.x specifications in YAML or JSON format.
// Swagger 2.x files are automatically converted to OpenAPI 3.x for uniform processing.
func Parse(path string) (*openapi3.T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}

	if isSwagger2(data) {
		return parseSwagger2(data)
	}

	return parseOpenAPI3(path)
}

func isSwagger2(data []byte) bool {
	content := string(data)
	return strings.Contains(content, `"swagger"`) ||
		strings.Contains(content, `swagger:`) ||
		strings.Contains(content, `"swagger":`)
}

func parseSwagger2(data []byte) (*openapi3.T, error) {
	var doc openapi2.T
	if err := doc.UnmarshalJSON(data); err != nil {
		return nil, fmt.Errorf("parsing Swagger 2.x spec: %w", err)
	}

	doc3, err := openapi2conv.ToV3(&doc)
	if err != nil {
		return nil, fmt.Errorf("converting Swagger 2.x to OpenAPI 3.x: %w", err)
	}

	return doc3, nil
}

func parseOpenAPI3(path string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("parsing OpenAPI 3.x spec: %w", err)
	}

	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("validating OpenAPI 3.x spec: %w", err)
	}

	return doc, nil
}
