package comparator_test

import (
	"testing"

	"github.com/yourorg/envlock/comparator"
	"github.com/yourorg/envlock/schema"
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

func TestCompare_Identical(t *testing.T) {
	s := makeSchema(ev("PORT", "port", "8080", "", true, false))
	result := comparator.Compare(s, s)
	if result.HasChanges() {
		t.Fatal("expected no changes for identical schemas")
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := makeSchema(ev("PORT", "port", "8080", "", true, false), ev("HOST", "host", "", "", true, false))
	right := makeSchema(ev("PORT", "port", "8080", "", true, false))
	result := comparator.Compare(left, right)
	if len(result.OnlyInLeft) != 1 || result.OnlyInLeft[0] != "HOST" {
		t.Fatalf("expected HOST only in left, got %v", result.OnlyInLeft)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := makeSchema(ev("PORT", "port", "8080", "", true, false))
	right := makeSchema(ev("PORT", "port", "8080", "", true, false), ev("DEBUG", "debug", "false", "", false, false))
	result := comparator.Compare(left, right)
	if len(result.OnlyInRight) != 1 || result.OnlyInRight[0] != "DEBUG" {
		t.Fatalf("expected DEBUG only in right, got %v", result.OnlyInRight)
	}
}

func TestCompare_FieldDiff(t *testing.T) {
	left := makeSchema(ev("PORT", "The port", "8080", "", true, false))
	right := makeSchema(ev("PORT", "Service port", "9090", "", true, false))
	result := comparator.Compare(left, right)
	if len(result.Differing) != 1 {
		t.Fatalf("expected 1 differing entry, got %d", len(result.Differing))
	}
	entry := result.Differing[0]
	if entry.Key != "PORT" {
		t.Fatalf("expected PORT, got %s", entry.Key)
	}
	if len(entry.Fields) != 2 {
		t.Fatalf("expected 2 field diffs, got %d", len(entry.Fields))
	}
}

func TestCompare_SensitiveFieldDiff(t *testing.T) {
	left := makeSchema(ev("SECRET", "a secret", "", "", true, false))
	right := makeSchema(ev("SECRET", "a secret", "", "", true, true))
	result := comparator.Compare(left, right)
	if len(result.Differing) != 1 {
		t.Fatal("expected differing entry for sensitive change")
	}
	fields := result.Differing[0].Fields
	if fields[0].Field != "sensitive" {
		t.Fatalf("expected sensitive field diff, got %s", fields[0].Field)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	s := makeSchema(
		ev("A", "desc a", "1", "", true, false),
		ev("B", "desc b", "", `\d+`, false, true),
	)
	result := comparator.Compare(s, s)
	if result.HasChanges() {
		t.Fatal("expected no changes")
	}
}
