package templater_test

import (
	"testing"

	"github.com/yourorg/envlock/schema"
	"github.com/yourorg/envlock/templater"
)

func makeSchema(vars []schema.VarEntry) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestRender_BasicSubstitution(t *testing.T) {
	s := makeSchema([]schema.VarEntry{
		{Key: "APP_HOST", Required: true},
		{Key: "APP_PORT", Required: true},
	})
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	result, err := templater.Render("http://${APP_HOST}:${APP_PORT}", env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "http://localhost:8080" {
		t.Errorf("expected 'http://localhost:8080', got %q", result.Output)
	}
	if len(result.Missing) != 0 {
		t.Errorf("expected no missing vars, got %v", result.Missing)
	}
}

func TestRender_DefaultApplied(t *testing.T) {
	s := makeSchema([]schema.VarEntry{
		{Key: "LOG_LEVEL", Required: false, Default: "info"},
	})
	result, err := templater.Render("level=${LOG_LEVEL}", map[string]string{}, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "level=info" {
		t.Errorf("expected 'level=info', got %q", result.Output)
	}
}

func TestRender_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.VarEntry{
		{Key: "DB_URL", Required: true},
	})
	result, err := templater.Render("dsn=${DB_URL}", map[string]string{}, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Missing) != 1 || result.Missing[0] != "DB_URL" {
		t.Errorf("expected DB_URL in missing, got %v", result.Missing)
	}
	if result.Output != "dsn=" {
		t.Errorf("expected 'dsn=', got %q", result.Output)
	}
}

func TestRender_NoPlaceholders(t *testing.T) {
	s := makeSchema([]schema.VarEntry{})
	result, err := templater.Render("static content", map[string]string{}, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "static content" {
		t.Errorf("expected 'static content', got %q", result.Output)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	s := makeSchema([]schema.VarEntry{})
	_, err := templater.Render("{{invalid", map[string]string{}, s)
	if err == nil {
		t.Error("expected error for invalid template")
	}
}
