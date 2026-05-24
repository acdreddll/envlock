package linter_test

import (
	"testing"

	"github.com/envlock/linter"
	"github.com/envlock/schema"
)

func TestLint_FullSchema(t *testing.T) {
	s := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_NAME", Description: "Application name", Required: true},
			{Key: "APP_PORT", Description: "Port to listen on", Default: "8080"},
			{Key: "DB_PASSWORD", Description: "Database password", Required: true, Sensitive: true, Pattern: `^.{12,}$`},
			{Key: "LOG_LEVEL", Description: "Logging verbosity", Default: "info"},
		},
	}

	issues := linter.Lint(s)
	if len(issues) != 0 {
		t.Errorf("expected no issues for well-annotated schema, got %d: %v", len(issues), issues)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	s := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_NAME"},
			{Key: "APP_NAME", Description: "Duplicate key"},
			{Key: "DB_PASS", Required: true, Default: "secret"},
			{Key: "API_TOKEN", Sensitive: true},
		},
	}

	issues := linter.Lint(s)

	kinds := map[string]int{}
	for _, iss := range issues {
		kinds[iss.Kind]++
	}

	if kinds["duplicate_key"] == 0 {
		t.Error("expected duplicate_key issue")
	}
	if kinds["required_with_default"] == 0 {
		t.Error("expected required_with_default issue")
	}
	if kinds["missing_description"] == 0 {
		t.Error("expected missing_description issue")
	}
}

func TestLint_EmptySchema(t *testing.T) {
	s := schema.Schema{}
	issues := linter.Lint(s)
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty schema, got %d", len(issues))
	}
}
