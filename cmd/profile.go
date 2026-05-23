package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/profiler"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var profileFormat string

func init() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile schema completeness and annotation coverage",
		RunE:  runProfile,
	}
	profileCmd.Flags().StringVar(&profileFormat, "format", "text", "Output format: text|json")
	rootCmd.AddCommand(profileCmd)
}

func runProfile(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	report := profiler.Profile(s)

	switch profileFormat {
	case "json":
		return printProfileJSON(report)
	default:
		printProfileText(report)
		return nil
	}
}

func printProfileText(r profiler.Report) {
	fmt.Println(profiler.Summary(r))
	if len(r.MissingDesc) > 0 {
		fmt.Println("\nMissing description:")
		for _, k := range r.MissingDesc {
			fmt.Printf("  - %s\n", k)
		}
	}
	if len(r.MissingDefault) > 0 {
		fmt.Println("\nOptional fields missing default:")
		for _, k := range r.MissingDefault {
			fmt.Printf("  - %s\n", k)
		}
	}
	if len(r.Fields) > 0 {
		fmt.Println("\nField scores:")
		for _, f := range r.Fields {
			fmt.Printf("  %-30s %3d%%\n", f.Key, f.CompletenessScore)
		}
	}
}

func printProfileJSON(r profiler.Report) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
