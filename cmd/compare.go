package cmd

import (
	"fmt"
	"os"

	"github.com/martyn/apidiff/internal/diff"
	"github.com/martyn/apidiff/internal/formatter"
	"github.com/martyn/apidiff/internal/parser"
	"github.com/spf13/cobra"
)

var (
	oldSpec        string
	newSpec        string
	format         string
	failOnBreaking bool
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two OpenAPI specifications",
	Long:  "Compare two OpenAPI/Swagger specifications and report differences including breaking changes.",
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().StringVar(&oldSpec, "old", "", "path to the old/original OpenAPI spec")
	compareCmd.Flags().StringVar(&newSpec, "new", "", "path to the new/updated OpenAPI spec")
	compareCmd.Flags().StringVar(&format, "format", "console", "output format: console, json, markdown, sarif")
	compareCmd.Flags().BoolVar(&failOnBreaking, "fail-on-breaking", false, "exit with non-zero code if breaking changes are found")

	_ = compareCmd.MarkFlagRequired("old")
	_ = compareCmd.MarkFlagRequired("new")

	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	oldDoc, err := parser.Parse(oldSpec)
	if err != nil {
		return fmt.Errorf("failed to parse old spec: %w", err)
	}

	newDoc, err := parser.Parse(newSpec)
	if err != nil {
		return fmt.Errorf("failed to parse new spec: %w", err)
	}

	changes := diff.Compare(oldDoc, newDoc)

	f, err := formatter.New(format)
	if err != nil {
		return err
	}

	output, err := f.Format(changes)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	fmt.Print(output)

	if failOnBreaking {
		summary := diff.Summarize(changes)
		if summary.Breaking > 0 {
			os.Exit(1)
		}
	}

	return nil
}
