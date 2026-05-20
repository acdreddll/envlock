package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/deprecator"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var deprecateCmd = &cobra.Command{
	Use:   "deprecate",
	Short: "List deprecated environment variables defined in the schema",
	RunE:  runDeprecate,
}

func init() {
	deprecateCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(deprecateCmd)
}

func runDeprecate(cmd *cobra.Command, args []string) error {
	schemaFile, _ := cmd.Flags().GetString("schema")
	format, _ := cmd.Flags().GetString("format")

	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	result := deprecator.Deprecate(s)

	switch format {
	case "json":
		return printDeprecateJSON(result)
	default:
		printDeprecateText(result)
		return nil
	}
}

func printDeprecateText(result deprecator.Result) {
	if !result.HasIssues() {
		fmt.Println("No deprecated variables found.")
		return
	}
	fmt.Printf("Found %d deprecated variable(s):\n\n", len(result.Findings))
	for _, f := range result.Findings {
		fmt.Printf("  [%s] %s\n", f.Key, f.Message)
	}
}

func printDeprecateJSON(result deprecator.Result) error {
	out := struct {
		Count    int                   `json:"count"`
		Findings []deprecator.Finding  `json:"findings"`
	}{
		Count:    len(result.Findings),
		Findings: result.Findings,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
