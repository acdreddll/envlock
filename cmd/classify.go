package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/envlock/classifier"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var classifyFormat string

func init() {
	classifyCmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify environment variables by category",
		Long:  "Assigns each variable a category: secret, config, feature_flag, infra, or unknown.",
		RunE:  runClassify,
	}
	classifyCmd.Flags().StringVarP(&classifyFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(classifyCmd)
}

func runClassify(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	results := classifier.Classify(s)

	switch classifyFormat {
	case "json":
		return printClassifyJSON(results)
	default:
		printClassifyText(results)
		return nil
	}
}

func printClassifyText(results []classifier.Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tCATEGORY\tREASON")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.Category, r.Reason)
	}
	w.Flush()
}

func printClassifyJSON(results []classifier.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
