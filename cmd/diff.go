package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envlock/differ"
	"github.com/yourorg/envlock/schema"
)

var diffCmd = &cobra.Command{
	Use:   "diff <base-schema> <next-schema>",
	Short: "Compare two envlock schema files and show differences",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	basePath := args[0]
	nextPath := args[1]

	baseSchema, err := schema.Load(basePath)
	if err != nil {
		return fmt.Errorf("loading base schema %q: %w", basePath, err)
	}

	nextSchema, err := schema.Load(nextPath)
	if err != nil {
		return fmt.Errorf("loading next schema %q: %w", nextPath, err)
	}

	d := differ.Compare(baseSchema, nextSchema)

	if !d.HasChanges() {
		fmt.Println("No differences found.")
		return nil
	}

	fmt.Printf("Schema diff: %s\n\n", d.Summary())

	for _, v := range d.Added {
		fmt.Printf("  [+] %s", v.Key)
		if v.Required {
			fmt.Print(" (required)")
		}
		if v.Description != "" {
			fmt.Printf(" — %s", v.Description)
		}
		fmt.Println()
	}

	for _, v := range d.Removed {
		fmt.Printf("  [-] %s\n", v.Key)
	}

	for _, c := range d.Changed {
		fmt.Printf("  [~] %s\n", c.Key)
		if c.From.Required != c.To.Required {
			fmt.Printf("       required: %v → %v\n", c.From.Required, c.To.Required)
		}
		if c.From.Default != c.To.Default {
			fmt.Printf("       default:  %q → %q\n", c.From.Default, c.To.Default)
		}
		if c.From.Pattern != c.To.Pattern {
			fmt.Printf("       pattern:  %q → %q\n", c.From.Pattern, c.To.Pattern)
		}
	}

	os.Exit(1)
	return nil
}
