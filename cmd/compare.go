package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envlock/comparator"
	"github.com/yourorg/envlock/schema"
)

var compareFormat string

func init() {
	compareCmd := &cobra.Command{
		Use:   "compare <left-schema> <right-schema>",
		Short: "Deep field-level comparison between two schema files",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	compareCmd.Flags().StringVarP(&compareFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	left, err := schema.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading left schema: %w", err)
	}
	right, err := schema.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading right schema: %w", err)
	}

	result := comparator.Compare(left, right)

	switch compareFormat {
	case "json":
		return printCompareJSON(result)
	default:
		printCompareText(result)
		return nil
	}
}

func printCompareText(r comparator.Result) {
	if !r.HasChanges() {
		fmt.Println("Schemas are identical.")
		return
	}
	for _, k := range r.OnlyInLeft {
		fmt.Printf("  only in left:  %s\n", k)
	}
	for _, k := range r.OnlyInRight {
		fmt.Printf("  only in right: %s\n", k)
	}
	for _, entry := range r.Differing {
		fmt.Printf("  changed: %s\n", entry.Key)
		for _, fd := range entry.Fields {
			fmt.Printf("    %s: %q -> %q\n", fd.Field, fd.Left, fd.Right)
		}
	}
}

func printCompareJSON(r comparator.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
