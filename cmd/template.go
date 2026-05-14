package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envlock/schema"
	"github.com/yourorg/envlock/templater"
)

var (
	templateSchemaFile string
	templateInput      string
	templateEnvFile    string
	failOnMissing      bool
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Render a template file with environment variable substitution",
		Long: `Reads a template file containing \${VAR} placeholders and substitutes
values from the environment (or a .env file), validated against the schema.`,
		RunE: runTemplate,
	}

	templateCmd.Flags().StringVarP(&templateSchemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	templateCmd.Flags().StringVarP(&templateInput, "input", "i", "", "Path to template file (required)")
	templateCmd.Flags().StringVarP(&templateEnvFile, "env-file", "e", "", "Optional .env file to load")
	templateCmd.Flags().BoolVar(&failOnMissing, "fail-on-missing", false, "Exit with error if required vars are missing")
	_ = templateCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(templateSchemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	tmplBytes, err := os.ReadFile(templateInput)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	env := parseProcessEnv(os.Environ())
	if templateEnvFile != "" {
		fileEnv, err := parseDotEnv(templateEnvFile)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}
		for k, v := range fileEnv {
			env[k] = v
		}
	}

	result, err := templater.Render(string(tmplBytes), env, s)
	if err != nil {
		return fmt.Errorf("template rendering failed: %w", err)
	}

	if failOnMissing && len(result.Missing) > 0 {
		for _, key := range result.Missing {
			fmt.Fprintf(os.Stderr, "missing required variable: %s\n", key)
		}
		return fmt.Errorf("%d required variable(s) missing", len(result.Missing))
	}

	fmt.Print(result.Output)
	return nil
}
