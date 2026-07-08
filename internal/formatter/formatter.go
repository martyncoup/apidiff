package formatter

import (
	"fmt"

	"github.com/martyn/apidiff/internal/model"
)

// Formatter formats a list of changes into an output string.
type Formatter interface {
	Format(changes []model.Change) (string, error)
}

// New returns a formatter for the given format name.
func New(format string) (Formatter, error) {
	switch format {
	case "console", "":
		return &TextFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "markdown", "md":
		return &MarkdownFormatter{}, nil
	case "sarif":
		return &SARIFFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: console, json, markdown, sarif)", format)
	}
}
