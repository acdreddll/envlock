package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlock/auditor"
	"github.com/user/envlock/schema"
)

var auditOutputFormat string

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit environment variables for security and hygiene issues",
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "path to schema file")
	auditCmd.Flags().StringVarP(&auditOutputFormat, "format", "f", "text", "output format: text or json")
	auditCmd.Flags().StringArrayVar(&dotEnvFiles, "env-file", nil, "path(s) to .env file(s) to load")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	env := collectEnv(dotEnvFiles)
	report := auditor.Audit(s, env)

	switch auditOutputFormat {
	case "json":
		return printAuditJSON(report)
	default:
		return printAuditText(report)
	}
}

func printAuditText(r auditor.Report) error {
	fmt.Println(r.Summary())
	for _, f := range r.Findings {
		fmt.Fprintf(os.Stdout, "  [%s] %s: %s\n", f.Severity, f.Key, f.Message)
	}
	if r.HasIssues() {
		os.Exit(1)
	}
	return nil
}

func printAuditJSON(r auditor.Report) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("encoding audit report: %w", err)
	}
	if r.HasIssues() {
		os.Exit(1)
	}
	return nil
}
