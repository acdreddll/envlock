package validator_test

import (
	"testing"

	"github.com/yourorg/envlock/schema"
	"github.com/yourorg/envlock/validator"
)

func TestSummary_AllPassed(t *testing.T) {
	sc := makeSchema(
		schema.EnvVar{Key: "FOO", Required: true},
		schema.EnvVar{Key: "BAR", Required: true},
	)
	report := validator.Validate(sc, map[string]string{"FOO": "a", "BAR": "b"})
	got := validator.Summary(report)
	want := "all 2 variable(s) passed validation"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSummary_WithErrors(t *testing.T) {
	sc := makeSchema(
		schema.EnvVar{Key: "FOO", Required: true},
		schema.EnvVar{Key: "BAR", Required: true},
		schema.EnvVar{Key: "BAZ", Required: true},
	)
	report := validator.Validate(sc, map[string]string{"FOO": "a"})
	got := validator.Summary(report)
	want := "1/3 variable(s) passed, 2 error(s) found"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSummary_EmptySchema(t *testing.T) {
	sc := schema.Schema{}
	report := validator.Validate(sc, map[string]string{})
	got := validator.Summary(report)
	want := "no variables defined in schema"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestErrorKeys_ReturnsFailedKeys(t *testing.T) {
	sc := makeSchema(
		schema.EnvVar{Key: "MISSING", Required: true},
		schema.EnvVar{Key: "PRESENT", Required: true},
	)
	report := validator.Validate(sc, map[string]string{"PRESENT": "yes"})
	keys := validator.ErrorKeys(report)
	if len(keys) != 1 || keys[0] != "MISSING" {
		t.Errorf("expected [MISSING], got %v", keys)
	}
}

func TestPassedKeys_ReturnsPassedKeys(t *testing.T) {
	sc := makeSchema(
		schema.EnvVar{Key: "A", Required: true},
		schema.EnvVar{Key: "B", Required: true},
	)
	report := validator.Validate(sc, map[string]string{"A": "1", "B": "2"})
	keys := validator.PassedKeys(report)
	if len(keys) != 2 {
		t.Errorf("expected 2 passed keys, got %v", keys)
	}
}
