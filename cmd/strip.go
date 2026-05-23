package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/envlock/schema"
	"github.com/your-org/envlock/stripper"
	"gopkg.in/yaml.v3"
)

var stripFields []string
var stripOutput string

func init() {
	stripCmd := &cobra.Command{
		Use:   "strip",
		Short: "Remove specified fields from schema variables",
		Long: `Strip removes one or more fields from every variable in the schema.
Useful for publishing a sanitised schema without sensitive metadata.

Valid fields: ` + strings.Join(stripper.ValidFields, ", "),
		RunE: runStrip,
	}

	stripCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "path to schema file")
	stripCmd.Flags().StringSliceVarP(&stripFields, "fields", "f", nil, "comma-separated list of fields to strip")
	stripCmd.Flags().StringVarP(&stripOutput, "output", "o", "yaml", "output format: yaml or json")
	_ = stripCmd.MarkFlagRequired("fields")

	rootCmd.AddCommand(stripCmd)
}

func runStrip(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	out, err := stripper.Strip(s, stripper.Options{Fields: stripFields})
	if err != nil {
		return fmt.Errorf("stripping fields: %w", err)
	}

	switch stripOutput {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(out)
	default:
		return fmt.Errorf("unsupported output format %q: use yaml or json", stripOutput)
	}
}
