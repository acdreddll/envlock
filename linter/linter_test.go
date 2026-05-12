package linter_test

import (
	"testing"

	"github.com/envlock/linter"
	"github.com/envlock/schema"
)

func makeSchema(vars []schema.EnvVar) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestLint_NoIssues(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "APP_ENV", Required: true, Description: "Application environment"},
		{Key: "PORT", Default: "8080", Description: "Server port"},
	})
	result := linter.Lint(s)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got: %s", result.Summary())
	}
	if result.HasErrors() {
		t.Error("expected HasErrors to be false")
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "APP_ENV", Required: true, Description: "env"},
		{Key: "APP_ENV", Required: false, Description: "duplicate"},
	})
	result := linter.Lint(s)
	if !result.HasErrors() {
		t.Error("expected an error for duplicate key")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "APP_ENV" && issue.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate key error for APP_ENV")
	}
}

func TestLint_RequiredWithDefault(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "SECRET", Required: true, Default: "fallback", Description: "a secret"},
	})
	result := linter.Lint(s)
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "SECRET" && issue.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for required+default combination")
	}
}

func TestLint_MissingDescription(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "DB_URL", Required: true},
	})
	result := linter.Lint(s)
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "DB_URL" && issue.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for missing description on DB_URL")
	}
}

func TestLint_Summary_NoIssues(t *testing.T) {
	s := makeSchema([]schema.EnvVar{})
	result := linter.Lint(s)
	if result.Summary() != "schema is valid: no issues found" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}
