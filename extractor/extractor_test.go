package extractor_test

import (
	"testing"

	"github.com/envlock/extractor"
	"github.com/envlock/schema"
)

func makeSchema() schema.Schema {
	return schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "DB_HOST", Required: true, Group: "database", Tags: []string{"infra"}},
			{Key: "DB_PASS", Required: true, Sensitive: true, Group: "database", Tags: []string{"infra", "secret"}},
			{Key: "APP_PORT", Required: false, Group: "app", Tags: []string{"network"}},
			{Key: "SECRET_KEY", Required: true, Sensitive: true, Group: "app", Tags: []string{"secret"}},
			{Key: "LOG_LEVEL", Required: false, Group: "app"},
		},
	}
}

func TestExtract_SensitiveOnly(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{SensitiveOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Fatalf("expected 2 sensitive vars, got %d", len(res.Vars))
	}
}

func TestExtract_RequiredOnly(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{RequiredOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 3 {
		t.Fatalf("expected 3 required vars, got %d", len(res.Vars))
	}
}

func TestExtract_ByGroup(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{Group: "database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Fatalf("expected 2 database vars, got %d", len(res.Vars))
	}
}

func TestExtract_ByTag(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{Tags: []string{"secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Fatalf("expected 2 secret-tagged vars, got %d", len(res.Vars))
	}
}

func TestExtract_ByPattern(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Fatalf("expected 2 DB_ vars, got %d", len(res.Vars))
	}
}

func TestExtract_InvalidPattern(t *testing.T) {
	_, err := extractor.Extract(makeSchema(), extractor.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestExtract_CombinedFilters(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{
		Group:         "app",
		SensitiveOnly: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 1 || res.Vars[0].Key != "SECRET_KEY" {
		t.Fatalf("expected only SECRET_KEY, got %+v", res.Vars)
	}
}

func TestExtract_NoFilters(t *testing.T) {
	res, err := extractor.Extract(makeSchema(), extractor.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 5 {
		t.Fatalf("expected all 5 vars, got %d", len(res.Vars))
	}
}
