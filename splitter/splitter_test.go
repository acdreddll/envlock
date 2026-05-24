package splitter_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/splitter"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, group string, required, sensitive bool, tags ...string) schema.EnvVar {
	return schema.EnvVar{
		Key:       key,
		Group:     group,
		Required:  required,
		Sensitive: sensitive,
		Tags:      tags,
	}
}

func TestSplit_ByGroup(t *testing.T) {
	s := makeSchema(
		ev("DB_HOST", "database", true, false),
		ev("DB_PASS", "database", true, true),
		ev("APP_PORT", "app", false, false),
	)
	r := splitter.ByGroup(s, "database")
	if len(r.Matched.Vars) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(r.Matched.Vars))
	}
	if len(r.Remainder.Vars) != 1 {
		t.Fatalf("expected 1 remainder, got %d", len(r.Remainder.Vars))
	}
}

func TestSplit_ByGroup_CaseInsensitive(t *testing.T) {
	s := makeSchema(
		ev("DB_HOST", "Database", true, false),
		ev("APP_PORT", "app", false, false),
	)
	r := splitter.ByGroup(s, "database")
	if len(r.Matched.Vars) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(r.Matched.Vars))
	}
}

func TestSplit_ByRequired(t *testing.T) {
	s := makeSchema(
		ev("A", "", true, false),
		ev("B", "", false, false),
		ev("C", "", true, false),
	)
	r := splitter.ByRequired(s)
	if len(r.Matched.Vars) != 2 {
		t.Fatalf("expected 2 required, got %d", len(r.Matched.Vars))
	}
	if len(r.Remainder.Vars) != 1 {
		t.Fatalf("expected 1 optional, got %d", len(r.Remainder.Vars))
	}
}

func TestSplit_BySensitive(t *testing.T) {
	s := makeSchema(
		ev("SECRET", "", true, true),
		ev("HOST", "", true, false),
	)
	r := splitter.BySensitive(s)
	if len(r.Matched.Vars) != 1 || r.Matched.Vars[0].Key != "SECRET" {
		t.Fatalf("unexpected sensitive result: %+v", r.Matched.Vars)
	}
	if len(r.Remainder.Vars) != 1 || r.Remainder.Vars[0].Key != "HOST" {
		t.Fatalf("unexpected remainder: %+v", r.Remainder.Vars)
	}
}

func TestSplit_ByTag(t *testing.T) {
	s := makeSchema(
		ev("A", "", false, false, "infra", "core"),
		ev("B", "", false, false, "core"),
		ev("C", "", false, false, "other"),
	)
	r := splitter.ByTag(s, "infra")
	if len(r.Matched.Vars) != 1 {
		t.Fatalf("expected 1 infra var, got %d", len(r.Matched.Vars))
	}
	if len(r.Remainder.Vars) != 2 {
		t.Fatalf("expected 2 remainder vars, got %d", len(r.Remainder.Vars))
	}
}

func TestSplit_EmptySchema(t *testing.T) {
	r := splitter.ByRequired(makeSchema())
	if len(r.Matched.Vars) != 0 || len(r.Remainder.Vars) != 0 {
		t.Fatal("expected empty results for empty schema")
	}
}

func TestSplit_Summary(t *testing.T) {
	s := makeSchema(
		ev("A", "", true, false),
		ev("B", "", false, false),
	)
	r := splitter.ByRequired(s)
	got := splitter.Summary(r)
	if got != "matched=1 remainder=1" {
		t.Fatalf("unexpected summary: %q", got)
	}
}
