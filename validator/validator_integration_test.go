package validator_test

import (
	"testing"

	"github.com/yourorg/envlock/schema"
	"github.com/yourorg/envlock/validator"
)

func TestValidate_FullSchema(t *testing.T) {
	sc := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_NAME", Required: true, Description: "Application name"},
			{Key: "PORT", Required: true, Pattern: `^\d+$`, Description: "Port number"},
			{Key: "DEBUG", Required: false, Default: "false", Description: "Debug mode"},
			{Key: "API_SECRET", Required: true, Sensitive: true, Description: "API secret key"},
		},
	}

	env := map[string]string{
		"APP_NAME":   "myapp",
		"PORT":       "8080",
		"API_SECRET": "supersecretvalue",
	}

	report := validator.Validate(sc, env)

	if report.HasErrors() {
		t.Fatalf("expected no errors, got: %v", report.Errors)
	}
	if len(report.Passed) != 4 {
		t.Errorf("expected 4 passed vars, got %d", len(report.Passed))
	}
}

func TestValidate_MultipleFailures(t *testing.T) {
	sc := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_NAME", Required: true, Description: "Application name"},
			{Key: "PORT", Required: true, Pattern: `^\d+$`, Description: "Port number"},
			{Key: "LOG_LEVEL", Required: true, Pattern: `^(debug|info|warn|error)$`, Description: "Log level"},
		},
	}

	env := map[string]string{
		"PORT":      "not-a-port",
		"LOG_LEVEL": "verbose",
	}

	report := validator.Validate(sc, env)

	if !report.HasErrors() {
		t.Fatal("expected errors but got none")
	}
	if len(report.Errors) != 3 {
		t.Errorf("expected 3 errors (missing + 2 pattern), got %d", len(report.Errors))
	}
}

func TestValidate_EmptyEnvAllOptional(t *testing.T) {
	sc := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "FEATURE_FLAG", Required: false, Default: "off"},
			{Key: "TIMEOUT", Required: false, Default: "30"},
		},
	}

	report := validator.Validate(sc, map[string]string{})

	if report.HasErrors() {
		t.Fatalf("expected no errors for all-optional schema, got: %v", report.Errors)
	}
}
