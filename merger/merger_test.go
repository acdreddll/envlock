package merger_test

import (
	"testing"

	"github.com/envlock/merger"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc, def, pattern string, required, sensitive bool) schema.EnvVar {
	return schema.EnvVar{
		Key:         key,
		Description: desc,
		Default:     def,
		Pattern:     pattern,
		Required:    required,
		Sensitive:   sensitive,
	}
}

func TestMerge_BaseOnly(t *testing.T) {
	base := makeSchema(ev("PORT", "port", "8080", "", false, false))
	res, err := merger.Merge(base, schema.Schema{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Schema.Vars) != 1 || res.Schema.Vars[0].Key != "PORT" {
		t.Errorf("expected PORT in merged schema")
	}
}

func TestMerge_OverrideAddsNewKey(t *testing.T) {
	base := makeSchema(ev("PORT", "port", "8080", "", false, false))
	override := makeSchema(ev("HOST", "hostname", "localhost", "", true, false))
	res, err := merger.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Schema.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(res.Schema.Vars))
	}
	if res.Schema.Vars[0].Key != "PORT" || res.Schema.Vars[1].Key != "HOST" {
		t.Errorf("unexpected order: %v", res.Schema.Vars)
	}
}

func TestMerge_OverrideWinsOnDescription(t *testing.T) {
	base := makeSchema(ev("PORT", "old desc", "", "", true, false))
	override := makeSchema(ev("PORT", "new desc", "", "", true, false))
	res, _ := merger.Merge(base, override)
	if res.Schema.Vars[0].Description != "new desc" {
		t.Errorf("expected override description, got %q", res.Schema.Vars[0].Description)
	}
}

func TestMerge_RequiredConflictWarning(t *testing.T) {
	base := makeSchema(ev("TOKEN", "token", "", "", true, true))
	override := makeSchema(ev("TOKEN", "", "", "", false, true))
	res, err := merger.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for required flag conflict")
	}
	if res.Schema.Vars[0].Required != false {
		t.Error("expected override required=false to win")
	}
}

func TestMerge_SensitiveConflictWarning(t *testing.T) {
	base := makeSchema(ev("SECRET", "sec", "", "", true, false))
	override := makeSchema(ev("SECRET", "", "", "", true, true))
	res, _ := merger.Merge(base, override)
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for sensitive flag conflict")
	}
	if !res.Schema.Vars[0].Sensitive {
		t.Error("expected override sensitive=true to win")
	}
}

func TestMerge_NoWarningsWhenCompatible(t *testing.T) {
	base := makeSchema(ev("DB_URL", "db", "", "", true, true))
	override := makeSchema(ev("DB_URL", "database url", "postgres://localhost/db", `^postgres://`, true, true))
	res, _ := merger.Merge(base, override)
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", res.Warnings)
	}
	if res.Schema.Vars[0].Pattern != `^postgres://` {
		t.Errorf("expected pattern from override")
	}
}
