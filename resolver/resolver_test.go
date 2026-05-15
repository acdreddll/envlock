package resolver_test

import (
	"testing"

	"github.com/yourorg/envlock/resolver"
	"github.com/yourorg/envlock/schema"
)

func makeSchema(vars []schema.EnvVar) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestResolve_FromEnv(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "APP_HOST", Required: true},
	})
	env := map[string]string{"APP_HOST": "localhost"}

	result := resolver.Resolve(s, env)

	if result.HasErrors() {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Resolutions[0].Source != "env" {
		t.Errorf("expected source 'env', got %q", result.Resolutions[0].Source)
	}
	if result.Resolutions[0].Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", result.Resolutions[0].Value)
	}
}

func TestResolve_FallsBackToDefault(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "APP_PORT", Default: "8080"},
	})
	env := map[string]string{}

	result := resolver.Resolve(s, env)

	if result.HasErrors() {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Resolutions[0].Source != "default" {
		t.Errorf("expected source 'default', got %q", result.Resolutions[0].Source)
	}
	if result.Resolutions[0].Value != "8080" {
		t.Errorf("expected value '8080', got %q", result.Resolutions[0].Value)
	}
}

func TestResolve_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "DB_PASSWORD", Required: true},
	})
	env := map[string]string{}

	result := resolver.Resolve(s, env)

	if !result.HasErrors() {
		t.Fatal("expected errors for missing required variable")
	}
	if result.Resolutions[0].Source != "missing" {
		t.Errorf("expected source 'missing', got %q", result.Resolutions[0].Source)
	}
}

func TestResolve_OptionalMissing(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "OPTIONAL_FLAG", Required: false},
	})
	env := map[string]string{}

	result := resolver.Resolve(s, env)

	if result.HasErrors() {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if result.Resolutions[0].Source != "missing" {
		t.Errorf("expected source 'missing', got %q", result.Resolutions[0].Source)
	}
}

func TestResolve_EnvOverridesDefault(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "LOG_LEVEL", Default: "info"},
	})
	env := map[string]string{"LOG_LEVEL": "debug"}

	result := resolver.Resolve(s, env)

	if result.Resolutions[0].Source != "env" {
		t.Errorf("expected source 'env', got %q", result.Resolutions[0].Source)
	}
	if result.Resolutions[0].Value != "debug" {
		t.Errorf("expected value 'debug', got %q", result.Resolutions[0].Value)
	}
}
