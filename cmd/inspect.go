package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/envlock/inspector"
	"github.com/your-org/envlock/schema"
)

var inspectKey string
var inspectFormat string

func init() {
	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect variables defined in the schema",
		RunE:  runInspect,
	}
	inspectCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Schema file")
	inspectCmd.Flags().StringVarP(&inspectKey, "key", "k", "", "Inspect a single key")
	inspectCmd.Flags().StringVarP(&inspectFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	if inspectKey != "" {
		info, ok := inspector.Find(s, inspectKey)
		if !ok {
			return fmt.Errorf("key %q not found in schema", inspectKey)
		}
		if inspectFormat == "json" {
			return json.NewEncoder(os.Stdout).Encode(info)
		}
		printVarInfo(info)
		return nil
	}

	r := inspector.Inspect(s)
	if inspectFormat == "json" {
		return json.NewEncoder(os.Stdout).Encode(r)
	}

	fmt.Printf("Total variables: %d\n\n", r.Total)
	for _, v := range r.Vars {
		printVarInfo(v)
		fmt.Println()
	}
	return nil
}

func printVarInfo(v inspector.VarInfo) {
	fmt.Printf("Key:         %s\n", v.Key)
	if v.Description != "" {
		fmt.Printf("Description: %s\n", v.Description)
	}
	fmt.Printf("Required:    %v\n", v.Required)
	fmt.Printf("Sensitive:   %v\n", v.Sensitive)
	if v.HasDefault {
		fmt.Printf("Default:     %s\n", v.Default)
	}
	if v.Group != "" {
		fmt.Printf("Group:       %s\n", v.Group)
	}
	if len(v.Tags) > 0 {
		fmt.Printf("Tags:        %s\n", strings.Join(v.Tags, ", "))
	}
	if v.Pattern != "" {
		fmt.Printf("Pattern:     %s\n", v.Pattern)
	}
	if v.Deprecated {
		fmt.Printf("Deprecated:  true")
		if v.RemoveBy != "" {
			fmt.Printf(" (remove by: %s)", v.RemoveBy)
		}
		fmt.Println()
	}
}
