package patcher_test

import (
	"testing"

	"github.com/envlock/patcher"
	"github.com/envlock/schema"
)

func makeSchema() schema.Schema {
	return schema.Schema{
		{Key: "DATABASE_URL", Description: "Primary DB", Required: true, Pattern: `^postgres://`},
		{Key: "LOG_LEVEL", Description: "Logging verbosity", Default: "info", Group: "logging"},
		{Key: "SECRET_KEY", Sensitive: true},
	}
}

func TestPatch_UpdateDescription(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "LOG_LEVEL", Field: "description", Value: "Updated desc"}}
	out, results, err := patcher.Patch(s, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Fatalf("expected op to be applied, reason: %s", results[0].Reason)
	}
	if results[0].OldVal != "Logging verbosity" {
		t.Errorf("expected old val 'Logging verbosity', got %q", results[0].OldVal)
	}
	for _, ev := range out {
		if ev.Key == "LOG_LEVEL" && ev.Description != "Updated desc" {
			t.Errorf("expected description 'Updated desc', got %q", ev.Description)
		}
	}
}

func TestPatch_UpdateDefault(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "LOG_LEVEL", Field: "default", Value: "debug"}}
	out, results, _ := patcher.Patch(s, ops)
	if !results[0].Applied {
		t.Fatal("expected applied")
	}
	for _, ev := range out {
		if ev.Key == "LOG_LEVEL" && ev.Default != "debug" {
			t.Errorf("expected default 'debug', got %q", ev.Default)
		}
	}
}

func TestPatch_KeyNotFound(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "NONEXISTENT", Field: "description", Value: "x"}}
	_, results, _ := patcher.Patch(s, ops)
	if results[0].Applied {
		t.Fatal("expected not applied")
	}
	if results[0].Reason == "" {
		t.Error("expected a reason")
	}
}

func TestPatch_InvalidField(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "LOG_LEVEL", Field: "required", Value: "true"}}
	_, results, _ := patcher.Patch(s, ops)
	if results[0].Applied {
		t.Fatal("expected not applied for unknown field")
	}
}

func TestPatch_InvalidPattern(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "DATABASE_URL", Field: "pattern", Value: "[invalid"}}
	_, results, _ := patcher.Patch(s, ops)
	if results[0].Applied {
		t.Fatal("expected not applied for invalid pattern")
	}
}

func TestPatch_DoesNotMutateOriginal(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{{Key: "LOG_LEVEL", Field: "group", Value: "ops"}}
	_, _, _ = patcher.Patch(s, ops)
	for _, ev := range s {
		if ev.Key == "LOG_LEVEL" && ev.Group != "logging" {
			t.Error("original schema was mutated")
		}
	}
}

func TestPatch_MultipleOps(t *testing.T) {
	s := makeSchema()
	ops := []patcher.PatchOp{
		{Key: "LOG_LEVEL", Field: "description", Value: "New desc"},
		{Key: "SECRET_KEY", Field: "group", Value: "security"},
	}
	_, results, _ := patcher.Patch(s, ops)
	for _, r := range results {
		if !r.Applied {
			t.Errorf("expected op on %s/%s to be applied", r.Key, r.Field)
		}
	}
}
