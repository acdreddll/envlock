package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envlock/schema"
	"github.com/your-org/envlock/trimmer"
)

var trimMode string
var trimOutput string

func init() {
	trimCmd := &cobra.Command{
		Use:   "trim",
		Short: "Remove unused or unannotated entries from the schema",
		RunE:  runTrim,
	}
	trimCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	trimCmd.Flags().StringVar(&trimMode, "mode", "bare", "Trim mode: bare | optional-no-default")
	trimCmd.Flags().StringVar(&trimOutput, "output", "text", "Output format: text | json")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	var mode trimmer.TrimMode
	switch trimMode {
	case "bare":
		mode = trimmer.TrimBare
	case "optional-no-default":
		mode = trimmer.TrimOptionalNoDefault
	default:
		return fmt.Errorf("unknown trim mode %q; use bare or optional-no-default", trimMode)
	}

	result := trimmer.Trim(s, mode)

	if trimOutput == "json" {
		return printTrimJSON(result)
	}
	return printTrimText(result)
}

func printTrimText(r trimmer.Result) error {
	if len(r.Removed) == 0 {
		fmt.Println("No entries to trim.")
		return nil
	}
	fmt.Printf("Removed %d entr(ies):\n", len(r.Removed))
	for _, key := range r.Removed {
		fmt.Printf("  - %s\n", key)
	}
	fmt.Printf("Kept %d entr(ies).\n", len(r.Kept))
	return nil
}

func printTrimJSON(r trimmer.Result) error {
	out := map[string]interface{}{
		"removed": r.Removed,
		"kept_count": len(r.Kept),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
