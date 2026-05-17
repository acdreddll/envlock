package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/schema"
	"github.com/envlock/summarizer"
	"github.com/spf13/cobra"
)

var summarizeFormat string

func init() {
	summarizeCmd := &cobra.Command{
		Use:   "summarize",
		Short: "Print a statistical summary of the schema",
		RunE:  runSummarize,
	}
	summarizeCmd.Flags().StringVarP(&summarizeFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(summarizeCmd)
}

func runSummarize(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}
	sum := summarizer.Summarize(s.Vars)
	switch summarizeFormat {
	case "json":
		return printSummarizeJSON(sum)
	default:
		printSummarizeText(sum)
		return nil
	}
}

func printSummarizeText(s summarizer.Summary) {
	fmt.Fprintf(os.Stdout, "Total:        %d\n", s.Total)
	fmt.Fprintf(os.Stdout, "Required:     %d\n", s.Required)
	fmt.Fprintf(os.Stdout, "Optional:     %d\n", s.Optional)
	fmt.Fprintf(os.Stdout, "Sensitive:    %d\n", s.Sensitive)
	fmt.Fprintf(os.Stdout, "With Default: %d\n", s.WithDefault)
	fmt.Fprintf(os.Stdout, "With Pattern: %d\n", s.WithPattern)
	if len(s.Groups) > 0 {
		fmt.Fprintln(os.Stdout, "\nGroups:")
		for _, g := range summarizer.SortedGroups(s) {
			fmt.Fprintf(os.Stdout, "  %-20s %d\n", g, s.Groups[g])
		}
	}
	if len(s.TagCounts) > 0 {
		fmt.Fprintln(os.Stdout, "\nTags:")
		for _, tag := range summarizer.SortedTags(s) {
			fmt.Fprintf(os.Stdout, "  %-20s %d\n", tag, s.TagCounts[tag])
		}
	}
}

func printSummarizeJSON(s summarizer.Summary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
