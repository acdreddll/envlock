package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlock/schema"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envlock.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp schema: %v", err)
	}
	return p
}

func TestLoad_ValidSchema(t *testing.T) {
	content := `
version: "1"
vars:
  DATABASE_URL:
    description: "Primary database connection string"
    required: true
  LOG_LEVEL:
    description: "Logging verbosity"
    required: false
    default: "info"
`
	p := writeTempSchema(t, content)
	s, err := schema.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Version != "1" {
		t.Errorf("expected version 1, got %q", s.Version)
	}
	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}
	if !s.Vars["DATABASE_URL"].Required {
		t.Error("DATABASE_URL should be required")
	}
	if s.Vars["LOG_LEVEL"].Default != "info" {
		t.Errorf("expected default 'info', got %q", s.Vars["LOG_LEVEL"].Default)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/envlock.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTempSchema(t, ": invalid: yaml: [")
	_, err := schema.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	p := writeTempSchema(t, "vars:\n  FOO:\n    required: true\n")
	s, err := schema.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Version != "1" {
		t.Errorf("expected default version '1', got %q", s.Version)
	}
}
