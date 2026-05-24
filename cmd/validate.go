package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envlock/schema"
	"github.com/yourorg/envlock/validator"
)

var (
	validateSchemaFile string
	validateDotEnvFile string
	validateFormat     string
	validateStrict     bool
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate environment variables against the schema",
		Long: `Validate reads environment variables from the current process (or a .env file)
and checks them against the envlock schema. It reports missing required variables,
pattern mismatches, and unknown keys when running in strict mode.`,
		RunE: runValidate,
	}

	validateCmd.Flags().StringVarP(&validateSchemaFile, "schema", "s", "envlock.yaml", "Path to the envlock schema file")
	validateCmd.Flags().StringVarP(&validateDotEnvFile, "env-file", "e", "", "Path to a .env file to validate (defaults to process environment)")
	validateCmd.Flags().StringVarP(&validateFormat, "format", "f", "text", "Output format: text or json")
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Fail on unknown keys not present in the schema")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	sch, err := schema.Load(validateSchemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	var env map[string]string
	if validateDotEnvFile != "" {
		env, err = parseDotEnv(validateDotEnvFile)
		if err != nil {
			return fmt.Errorf("reading env file: %w", err)
		}
	} else {
		env = parseProcessEnv(os.Environ())
	}

	report := validator.Validate(sch, env)

	if validateStrict {
		schemaKeys := make(map[string]bool, len(sch.Vars))
		for _, v := range sch.Vars {
			schemaKeys[v.Key] = true
		}
		for k := range env {
			if !schemaKeys[k] {
				report = append(report, validator.Result{
					Key:    k,
					OK:     false,
					Reason: "unknown key not present in schema (strict mode)",
				})
			}
		}
	}

	switch validateFormat {
	case "json":
		printValidateJSON(report)
	default:
		printValidateText(report)
	}

	sum := validator.Summary(report)
	if sum.Failed > 0 {
		os.Exit(1)
	}
	return nil
}

func printValidateText(report []validator.Result) {
	passed := 0
	for _, r := range report {
		if r.OK {
			passed++
			fmt.Printf("  ✓  %s\n", r.Key)
		} else {
			fmt.Printf("  ✗  %s — %s\n", r.Key, r.Reason)
		}
	}
	sum := validator.Summary(report)
	fmt.Printf("\n%d passed, %d failed\n", sum.Passed, sum.Failed)
}

func printValidateJSON(report []validator.Result) {
	type jsonOut struct {
		Results []validator.Result `json:"results"`
		Passed  int                `json:"passed"`
		Failed  int                `json:"failed"`
	}
	sum := validator.Summary(report)
	out := jsonOut{
		Results: report,
		Passed:  sum.Passed,
		Failed:  sum.Failed,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
