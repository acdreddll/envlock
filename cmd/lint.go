package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/linter"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var lintFormat string

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Check schema for style and correctness issues",
		RunE:  runLint,
	}
	lintCmd.Flags().StringVarP(&lintFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	issues := linter.Lint(s)
	summary := linter.Summarize(issues)

	switch lintFormat {
	case "json":
		return printLintJSON(issues, summary)
	default:
		printLintText(issues, summary)
		if summary.HasIssues() {
			os.Exit(1)
		}
		return nil
	}
}

func printLintText(issues []linter.Issue, summary linter.Summary) {
	if !summary.HasIssues() {
		fmt.Println("No lint issues found.")
		return
	}
	for _, iss := range issues {
		fmt.Printf("[%s] %s: %s\n", iss.Kind, iss.Key, iss.Message)
	}
	fmt.Printf("\n%d issue(s) found across %d key(s).\n", summary.Total, len(summary.Affected))
}

func printLintJSON(issues []linter.Issue, summary linter.Summary) error {
	out := struct {
		Issues  []linter.Issue `json:"issues"`
		Total   int            `json:"total"`
		ByKind  map[string]int `json:"by_kind"`
		Affected []string      `json:"affected_keys"`
	}{
		Issues:   issues,
		Total:    summary.Total,
		ByKind:   summary.ByKind,
		Affected: summary.Affected,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	if summary.HasIssues() {
		os.Exit(1)
	}
	return nil
}
