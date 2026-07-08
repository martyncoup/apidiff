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

type jsonOutputWithVersion struct {
	Summary            jsonSummary    `json:"summary"`
	RecommendedVersion string         `json:"recommended_version,omitempty"`
	Changes            []model.Change `json:"changes"`
}

func (f *JSONFormatter) Format(changes []model.Change, opts Options) (string, error) {
	s := diff.Summarize(changes)

	output := jsonOutputWithVersion{
		Summary: jsonSummary{
			Added:         s.Added,
			Removed:       s.Removed,
			SchemaChanges: s.SchemaChanges,
			Breaking:      s.Breaking,
		},
		Changes: changes,
	}

	if opts.RecommendVersion {
		output.RecommendedVersion = string(diff.RecommendVersion(changes))
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
