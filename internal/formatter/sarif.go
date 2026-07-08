package formatter

import (
	"encoding/json"

	"github.com/martyn/apidiff/internal/model"
)

// SARIFFormatter outputs changes in SARIF v2.1.0 format for CI integration.
type SARIFFormatter struct{}

type sarifLog struct {
	Schema  string      `json:"$schema"`
	Version string      `json:"version"`
	Runs    []sarifRun  `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string          `json:"id"`
	ShortDescription sarifMessage    `json:"shortDescription"`
	DefaultConfig    sarifRuleConfig `json:"defaultConfiguration"`
}

type sarifRuleConfig struct {
	Level string `json:"level"`
}

type sarifResult struct {
	RuleID  string       `json:"ruleId"`
	Level   string       `json:"level"`
	Message sarifMessage `json:"message"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

func (f *SARIFFormatter) Format(changes []model.Change) (string, error) {
	rules := buildRules()
	results := make([]sarifResult, 0, len(changes))

	for _, c := range changes {
		results = append(results, sarifResult{
			RuleID:  string(c.Type),
			Level:   sarifLevel(c.Severity),
			Message: sarifMessage{Text: c.Description},
		})
	}

	log := sarifLog{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "apidiff",
						Version:        "0.1.0",
						InformationURI: "https://github.com/martyn/apidiff",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func buildRules() []sarifRule {
	return []sarifRule{
		{ID: string(model.EndpointAdded), ShortDescription: sarifMessage{Text: "An endpoint was added"}, DefaultConfig: sarifRuleConfig{Level: "note"}},
		{ID: string(model.EndpointRemoved), ShortDescription: sarifMessage{Text: "An endpoint was removed"}, DefaultConfig: sarifRuleConfig{Level: "error"}},
		{ID: string(model.PropertyAdded), ShortDescription: sarifMessage{Text: "A property was added"}, DefaultConfig: sarifRuleConfig{Level: "note"}},
		{ID: string(model.PropertyRemoved), ShortDescription: sarifMessage{Text: "A property was removed"}, DefaultConfig: sarifRuleConfig{Level: "error"}},
		{ID: string(model.PropertyTypeChanged), ShortDescription: sarifMessage{Text: "A property type was changed"}, DefaultConfig: sarifRuleConfig{Level: "error"}},
		{ID: string(model.SchemaChanged), ShortDescription: sarifMessage{Text: "A schema was changed"}, DefaultConfig: sarifRuleConfig{Level: "warning"}},
	}
}

func sarifLevel(severity model.Severity) string {
	switch severity {
	case model.SeverityBreaking:
		return "error"
	case model.SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}
