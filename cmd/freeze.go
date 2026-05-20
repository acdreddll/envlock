package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlock/freezer"
	"github.com/user/envlock/schema"
)

var (
	freezeSchemaFile string
	freezeEnvFile    string
	freezeFormat     string
	freezeRedact     bool
)

func init() {
	freezeCmd := &cobra.Command{
		Use:   "freeze",
		Short: "Capture a locked snapshot of resolved environment variables",
		RunE:  runFreeze,
	}
	freezeCmd.Flags().StringVarP(&freezeSchemaFile, "schema", "s", "envlock.yaml", "Schema file path")
	freezeCmd.Flags().StringVarP(&freezeEnvFile, "env-file", "e", "", "Optional .env file to load")
	freezeCmd.Flags().StringVarP(&freezeFormat, "format", "f", "text", "Output format: text or json")
	freezeCmd.Flags().BoolVar(&freezeRedact, "redact", false, "Mask sensitive values in output")
	rootCmd.AddCommand(freezeCmd)
}

func runFreeze(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(freezeSchemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	env := collectEnv(freezeEnvFile)

	snap, err := freezer.Freeze(s, env, freezeSchemaFile)
	if err != nil {
		return fmt.Errorf("freeze failed: %w", err)
	}

	if freezeRedact {
		snap = snap.Redacted()
	}

	switch freezeFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(snap)
	default:
		fmt.Printf("Snapshot created at: %s\n", snap.CreatedAt.Format("2006-01-02T15:04:05Z"))
		fmt.Printf("Schema: %s\n", snap.SchemaFile)
		fmt.Printf("Entries: %d\n\n", len(snap.Entries))
		for _, e := range snap.Entries {
			defaultMark := ""
			if e.FromDefault {
				defaultMark = " (default)"
			}
			sensitiveMark := ""
			if e.Sensitive {
				sensitiveMark = " [sensitive]"
			}
			fmt.Printf("  %-30s = %s%s%s\n", e.Key, e.Value, defaultMark, sensitiveMark)
		}
	}
	return nil
}
