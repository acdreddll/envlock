package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envlock/pruner"
	"github.com/your-org/envlock/schema"
	"gopkg.in/yaml.v3"
)

var (
	pruneDeprecated        bool
	pruneOptionalNoDefault bool
	pruneKeys              []string
	pruneOutput            string
	pruneFormat            string
)

func init() {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove unused or unwanted variables from the schema",
		RunE:  runPrune,
	}
	pruneCmd.Flags().BoolVar(&pruneDeprecated, "deprecated", false, "Remove deprecated entries")
	pruneCmd.Flags().BoolVar(&pruneOptionalNoDefault, "bare", false, "Remove optional vars with no default and no description")
	pruneCmd.Flags().StringArrayVar(&pruneKeys, "key", nil, "Explicitly remove a key (repeatable)")
	pruneCmd.Flags().StringVarP(&pruneOutput, "out", "o", "", "Write updated schema to file (default: stdout)")
	pruneCmd.Flags().StringVar(&pruneFormat, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	opts := pruner.Options{
		RemoveDeprecated:        pruneDeprecated,
		RemoveOptionalNoDefault: pruneOptionalNoDefault,
		Keys:                    pruneKeys,
	}

	res := pruner.Prune(s, opts)

	if pruneFormat == "json" {
		return json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"kept":    res.Kept,
			"removed": res.Removed,
		})
	}

	fmt.Printf("Pruned %d variable(s), kept %d.\n", len(res.Removed), len(res.Kept))
	for _, v := range res.Removed {
		fmt.Printf("  - %s\n", v.Key)
	}

	if pruneOutput != "" {
		updated := schema.Schema{Vars: res.Kept}
		f, err := os.Create(pruneOutput)
		if err != nil {
			return fmt.Errorf("open output file: %w", err)
		}
		defer f.Close()
		if err := yaml.NewEncoder(f).Encode(updated); err != nil {
			return fmt.Errorf("write schema: %w", err)
		}
		fmt.Printf("Updated schema written to %s\n", pruneOutput)
	}

	return nil
}
