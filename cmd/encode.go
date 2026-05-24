package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/encoder"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var encodeFormat string
var encodeSensitiveOnly bool

func init() {
	encodeCmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode environment variable values using base64 or hex",
		RunE:  runEncode,
	}
	encodeCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	encodeCmd.Flags().StringVarP(&encodeFormat, "format", "f", "base64", "Encoding format: base64, hex, none")
	encodeCmd.Flags().BoolVar(&encodeSensitiveOnly, "sensitive-only", false, "Only encode sensitive variables")
	encodeCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(encodeCmd)
}

func runEncode(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	env := collectEnv(cmd)

	results, err := encoder.Encode(s, env, encoder.Options{
		Format:        encoder.Format(encodeFormat),
		SensitiveOnly: encodeSensitiveOnly,
	})
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	if outputFormat == "json" {
		return printEncodeJSON(results)
	}
	return printEncodeText(results)
}

func printEncodeText(results []encoder.Result) error {
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(os.Stdout, "%-30s  [skipped]\n", r.Key)
			continue
		}
		fmt.Fprintf(os.Stdout, "%-30s  %s\n", r.Key, r.Encoded)
	}
	return nil
}

func printEncodeJSON(results []encoder.Result) error {
	return json.NewEncoder(os.Stdout).Encode(results)
}
