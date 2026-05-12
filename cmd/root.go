package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envlock/schema"
	"github.com/envlock/validator"
)

var (
	schemaFile string
	envFile    string
)

var rootCmd = &cobra.Command{
	Use:   "envlock",
	Short: "Validate environment variables against a schema contract",
	Long: `envlock validates that required environment variables are present
and match the expected patterns defined in your schema file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := schema.Load(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to load schema: %w", err)
		}

		env := collectEnv(envFile)

		report := validator.Validate(s, env)
		report.Print(os.Stdout)

		if !report.OK() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "path to the schema file")
	rootCmd.Flags().StringVarP(&envFile, "env-file", "e", "", "optional path to a .env file to validate against")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// collectEnv reads environment variables from the process environment.
// If an envFile path is provided, it parses and merges those values,
// with the env file values taking precedence over process environment values.
func collectEnv(envFile string) map[string]string {
	env := parseProcessEnv()

	if envFile != "" {
		parsed, err := parseDotEnv(envFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read env file %q: %v\n", envFile, err)
		} else {
			for k, v := range parsed {
				env[k] = v
			}
		}
	}

	return env
}

// parseProcessEnv returns a map of key-value pairs from the current process
// environment by splitting each entry on the first '=' character.
func parseProcessEnv() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return env
}
