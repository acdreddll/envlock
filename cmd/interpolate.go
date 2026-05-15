package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/interpolator"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var interpolateFormat string

func init() {
	interpolateCmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Expand cross-variable references in the resolved environment",
		RunE:  runInterpolate,
	}
	interpolateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	interpolateCmd.Flags().StringVar(&dotEnvFile, "dotenv", "", "Optional .env file to load")
	interpolateCmd.Flags().StringVar(&interpolateFormat, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	env := collectEnv(dotEnvFile)

	results, err := interpolator.Interpolate(s, env)
	if err != nil {
		return fmt.Errorf("interpolation failed: %w", err)
	}

	switch interpolateFormat {
	case "json":
		return printInterpolateJSON(results)
	default:
		printInterpolateText(results)
		return nil
	}
}

func printInterpolateText(results []interpolator.Result) {
	expanded := 0
	for _, r := range results {
		if r.Expanded {
			fmt.Printf("  ~ %-30s %q -> %q\n", r.Key, r.Original, r.Resolved)
			expanded++
		}
	}
	if expanded == 0 {
		fmt.Println("No variables required interpolation.")
	} else {
		fmt.Printf("\n%d variable(s) expanded.\n", expanded)
	}
}

func printInterpolateJSON(results []interpolator.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
