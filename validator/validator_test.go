package validator_test

import (
	"os"
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/validator"
)

func makeSchema(vars []schema.Var) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestValidate_AllPresent(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_PORT")

	s := makeSchema([]schema.Var{
		{Name: "APP_PORT", Required: true},
	})

	report, err := validator.Validate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Valid {
		t.Errorf("expected report to be valid")
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	os.Unsetenv("SECRET_KEY")

	s := makeSchema([]schema.Var{
		{Name: "SECRET_KEY", Required: true},
	})

	report, err := validator.Validate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Valid {
		t.Errorf("expected report to be invalid")
	}
	if len(report.Results) != 1 || !report.Results[0].Missing {
		t.Errorf("expected one missing result, got %+v", report.Results)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	os.Setenv("APP_ENV", "staging123")
	defer os.Unsetenv("APP_ENV")

	s := makeSchema([]schema.Var{
		{Name: "APP_ENV", Required: true, Pattern: "^(development|staging|production)$"},
	})

	report, err := validator.Validate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Valid {
		t.Errorf("expected report to be invalid due to pattern mismatch")
	}
	if len(report.Results) != 1 || !report.Results[0].Invalid {
		t.Errorf("expected one invalid result, got %+v", report.Results)
	}
}

func TestValidate_DefaultSkipsRequired(t *testing.T) {
	os.Unsetenv("LOG_LEVEL")

	s := makeSchema([]schema.Var{
		{Name: "LOG_LEVEL", Required: true, Default: "info"},
	})

	report, err := validator.Validate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Valid {
		t.Errorf("expected report to be valid when default is provided")
	}
}

func TestValidate_NilSchema(t *testing.T) {
	_, err := validator.Validate(nil)
	if err == nil {
		t.Errorf("expected error for nil schema")
	}
}
