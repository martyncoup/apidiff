package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/martyn/apidiff/internal/diff"
	"github.com/martyn/apidiff/internal/model"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	successIcon = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Render("✓")

	errorIcon = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Render("✗")

	breakingHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	propertyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			PaddingLeft(2)

	summaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)

// TextFormatter outputs human-readable styled console output using lipgloss.
type TextFormatter struct{}

func (f *TextFormatter) Format(changes []model.Change, opts Options) (string, error) {
	var b strings.Builder

	summary := diff.Summarize(changes)

	b.WriteString(titleStyle.Render("Comparing OpenAPI specifications..."))
	b.WriteString("\n\n")

	if summary.Added > 0 {
		b.WriteString(fmt.Sprintf("%s %s\n", successIcon, summaryStyle.Render(fmt.Sprintf("%d endpoint%s added", summary.Added, plural(summary.Added)))))
	}
	if summary.Removed > 0 {
		b.WriteString(fmt.Sprintf("%s %s\n", successIcon, summaryStyle.Render(fmt.Sprintf("%d endpoint%s removed", summary.Removed, plural(summary.Removed)))))
	}
	if summary.SchemaChanges > 0 {
		b.WriteString(fmt.Sprintf("%s %s\n", successIcon, summaryStyle.Render(fmt.Sprintf("%d schema change%s", summary.SchemaChanges, plural(summary.SchemaChanges)))))
	}

	// Breaking changes section
	breaking := filterBreaking(changes)
	if len(breaking) > 0 {
		b.WriteString("\n")
		b.WriteString(breakingHeader.Render("Breaking Changes"))
		b.WriteString("\n\n")
		for _, c := range breaking {
			b.WriteString(fmt.Sprintf("%s %s\n", errorIcon, c.Description))
			if c.Property != "" {
				b.WriteString(propertyStyle.Render(c.Property))
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
	}

	// Non-breaking changes section
	nonBreaking := filterNonBreaking(changes)
	if len(nonBreaking) > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %d non-breaking change%s omitted (use --format json for full details)", len(nonBreaking), plural(len(nonBreaking)))))
		b.WriteString("\n")
	}

	// Version recommendation
	if opts.RecommendVersion {
		bump := diff.RecommendVersion(changes)
		versionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%s Recommended version bump: %s\n", successIcon, versionStyle.Render(strings.ToUpper(string(bump)))))
	}

	return b.String(), nil
}

func filterBreaking(changes []model.Change) []model.Change {
	var breaking []model.Change
	for _, c := range changes {
		if c.IsBreaking() {
			breaking = append(breaking, c)
		}
	}
	return breaking
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
