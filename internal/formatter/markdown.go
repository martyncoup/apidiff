package formatter

import (
	"fmt"
	"strings"

	"github.com/martyn/apidiff/internal/diff"
	"github.com/martyn/apidiff/internal/model"
)

// MarkdownFormatter outputs changes as a Markdown document.
type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(changes []model.Change) (string, error) {
	var b strings.Builder

	summary := diff.Summarize(changes)

	b.WriteString("# API Diff Report\n\n")
	b.WriteString("## Summary\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Count |\n"))
	b.WriteString(fmt.Sprintf("|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Endpoints added | %d |\n", summary.Added))
	b.WriteString(fmt.Sprintf("| Endpoints removed | %d |\n", summary.Removed))
	b.WriteString(fmt.Sprintf("| Schema changes | %d |\n", summary.SchemaChanges))
	b.WriteString(fmt.Sprintf("| Breaking changes | %d |\n", summary.Breaking))

	breaking := filterBreaking(changes)
	if len(breaking) > 0 {
		b.WriteString("\n## Breaking Changes\n\n")
		b.WriteString("| Type | Path | Description |\n")
		b.WriteString("|------|------|-------------|\n")
		for _, c := range breaking {
			property := c.Property
			if property == "" {
				property = "-"
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", c.Type, property, c.Description))
		}
	}

	nonBreaking := filterNonBreaking(changes)
	if len(nonBreaking) > 0 {
		b.WriteString("\n## Other Changes\n\n")
		b.WriteString("| Type | Path | Description |\n")
		b.WriteString("|------|------|-------------|\n")
		for _, c := range nonBreaking {
			property := c.Property
			if property == "" {
				property = c.Path
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", c.Type, property, c.Description))
		}
	}

	return b.String(), nil
}

func filterNonBreaking(changes []model.Change) []model.Change {
	var result []model.Change
	for _, c := range changes {
		if !c.IsBreaking() {
			result = append(result, c)
		}
	}
	return result
}
