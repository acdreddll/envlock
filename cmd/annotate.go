package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/envlock/annotator"
	"github.com/your-org/envlock/schema"
	"gopkg.in/yaml.v3"
)

var (
	annotateKey         string
	annotateDescription string
	annotateGroup       string
	anotateTags         string
	annotateExample     string
	annotateOverwrite   bool
	annotateOutput      string
)

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Add or update metadata on a schema variable",
		RunE:  runAnnotate,
	}
	annotateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "schema file")
	annotateCmd.Flags().StringVarP(&annotateKey, "key", "k", "", "variable key to annotate (required)")
	annotateCmd.Flags().StringVar(&annotateDescription, "description", "", "description to set")
	annotateCmd.Flags().StringVar(&annotateGroup, "group", "", "group to set")
	annotateCmd.Flags().StringVar(&anotateTags, "tags", "", "comma-separated tags to add")
	annotateCmd.Flags().StringVar(&annotateExample, "example", "", "example value to set")
	annotateCmd.Flags().BoolVar(&annotateOverwrite, "overwrite", false, "overwrite existing values")
	annotateCmd.Flags().StringVarP(&annotateOutput, "output", "o", "yaml", "output format: yaml|json")
	_ = annotateCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	var tags []string
	if anotateTags != "" {
		for _, t := range strings.Split(anotateTags, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	opts := annotator.Options{
		Description: annotateDescription,
		Group:       annotateGroup,
		Tags:        tags,
		Example:     annotateExample,
		Overwrite:   annotateOverwrite,
	}

	updated, res, err := annotator.Annotate(s, annotateKey, opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "key=%s applied=%v skipped=%v\n", res.Key, res.Applied, res.Skipped)

	switch annotateOutput {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(updated)
	default:
		return yaml.NewEncoder(os.Stdout).Encode(updated)
	}
}
