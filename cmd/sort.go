package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/schema"
	"github.com/envlock/sorter"
	"github.com/spf13/cobra"
)

var sortOrder string
var sortOutputFormat string

func init() {
	sortCmd := &cobra.Command{
		Use:   "sort",
		Short: "Sort schema variables by key, required status, or description",
		RunE:  runSort,
	}

	sortCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	sortCmd.Flags().StringVarP(&sortOrder, "order", "o", "key", "Sort order: key, required, description")
	sortCmd.Flags().StringVarP(&sortOutputFormat, "format", "f", "text", "Output format: text, json")

	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	result := sorter.Sort(s, sorter.SortOrder(sortOrder))

	switch sortOutputFormat {
	case "json":
		return printSortJSON(result)
	default:
		printSortText(result)
		return nil
	}
}

func printSortText(result sorter.Result) {
	fmt.Printf("Sorted by: %s\n\n", result.Order)
	for _, ev := range result.Vars {
		required := "optional"
		if ev.Required {
			required = "required"
		}
		desc := ev.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  %-30s [%s] %s\n", ev.Key, required, desc)
	}
}

func printSortJSON(result sorter.Result) error {
	out := struct {
		Order string             `json:"order"`
		Vars  []schema.EnvVar    `json:"vars"`
	}{
		Order: string(result.Order),
		Vars:  result.Vars,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
