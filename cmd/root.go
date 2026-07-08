package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apidiff",
	Short: "Compare OpenAPI specifications and detect breaking changes",
	Long:  "apidiff is a CLI tool that compares two OpenAPI/Swagger specifications and reports added, removed, and changed endpoints and schemas.",
}

func Execute() error {
	return rootCmd.Execute()
}
