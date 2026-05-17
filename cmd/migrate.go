package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/migrator"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var migratePlanFile string
var migrateOutput string
var migrateFormat string

func init() {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Apply a migration plan to rename environment variable keys",
		RunE:  runMigrate,
	}
	migrateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	migrateCmd.Flags().StringVarP(&migratePlanFile, "plan", "p", "", "Path to migration plan YAML file (required)")
	migrateCmd.Flags().StringVarP(&migrateOutput, "output", "o", "", "Write updated schema to file (default: stdout)")
	migrateCmd.Flags().StringVarP(&migrateFormat, "format", "f", "text", "Output format: text or json")
	_ = migrateCmd.MarkFlagRequired("plan")
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	planData, err := os.ReadFile(migratePlanFile)
	if err != nil {
		return fmt.Errorf("reading plan file: %w", err)
	}

	var plan migrator.Plan
	if err := yaml.Unmarshal(planData, &plan); err != nil {
		return fmt.Errorf("parsing plan: %w", err)
	}

	updated, result, err := migrator.Apply(s, plan)
	if err != nil {
		return fmt.Errorf("applying migration: %w", err)
	}

	if migrateFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]interface{}{
			"applied":   result.Applied,
			"skipped":   result.Skipped,
			"conflicts": result.Conflicts,
		})
	}

	fmt.Printf("Applied:  %d\n", len(result.Applied))
	fmt.Printf("Skipped:  %d\n", len(result.Skipped))
	fmt.Printf("Conflicts:%d\n", len(result.Conflicts))
	for _, c := range result.Conflicts {
		fmt.Fprintf(os.Stderr, "  conflict: %s\n", c)
	}

	out, err := yaml.Marshal(updated)
	if err != nil {
		return fmt.Errorf("serializing updated schema: %w", err)
	}

	if migrateOutput != "" {
		return os.WriteFile(migrateOutput, out, 0o644)
	}
	fmt.Print(string(out))
	return nil
}
