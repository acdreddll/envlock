package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlock/schema"
	"github.com/user/envlock/scorer"
)

var scoreFormat string

func init() {
	scoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Score the quality of an envlock schema",
		RunE:  runScore,
	}
	scoreCmd.Flags().StringVarP(&scoreFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	result := scorer.Score(s)

	switch scoreFormat {
	case "json":
		return printScoreJSON(result)
	default:
		printScoreText(result)
		return nil
	}
}

func printScoreText(r scorer.Result) {
	fmt.Printf("Schema Quality Score: %d/100 (%.1f%%)\n", r.Score, r.Percentage)
	fmt.Printf("  Total vars:   %d\n", r.Total)
	fmt.Printf("  Descriptions: %d\n", r.Breakdown["description"])
	fmt.Printf("  Patterns:     %d\n", r.Breakdown["pattern"])
	fmt.Printf("  Defaults:     %d\n", r.Breakdown["default"])
	fmt.Printf("  Groups:       %d\n", r.Breakdown["group"])
}

func printScoreJSON(r scorer.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
