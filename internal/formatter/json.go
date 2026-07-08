package formatter

import (
	"encoding/json"

	"github.com/martyn/apidiff/internal/diff"
	"github.com/martyn/apidiff/internal/model"
)

// JSONFormatter outputs changes as a JSON document.
type JSONFormatter struct{}

type jsonOutput struct {
	Summary jsonSummary    `json:"summary"`
	Changes []model.Change `json:"changes"`
}

type jsonSummary struct {
	Added         int `json:"added"`
	Removed       int `json:"removed"`
	SchemaChanges int `json:"schema_changes"`
	Breaking      int `json:"breaking"`
}

func (f *JSONFormatter) Format(changes []model.Change) (string, error) {
	s := diff.Summarize(changes)

	output := jsonOutput{
		Summary: jsonSummary{
			Added:         s.Added,
			Removed:       s.Removed,
			SchemaChanges: s.SchemaChanges,
			Breaking:      s.Breaking,
		},
		Changes: changes,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
