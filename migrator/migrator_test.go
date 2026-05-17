package migrator_test

import (
	"testing"

	"github.com/envlock/migrator"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc string, required bool) schema.EnvVar {
	return schema.EnvVar{Key: key, Description: desc, Required: required}
}

func TestApply_RenamesKey(t *testing.T) {
	s := makeSchema(ev("OLD_KEY", "old", true))
	plan := migrator.Plan{
		Migrations: []migrator.Migration{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}},
	}
	out, result, err := migrator.Apply(s, plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 1 {
		t.Fatalf("expected 1 applied, got %d", len(result.Applied))
	}
	if len(out.Vars) != 1 || out.Vars[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in output, got %+v", out.Vars)
	}
}

func TestApply_SkipsMissingOldKey(t *testing.T) {
	s := makeSchema(ev("EXISTING", "x", false))
	plan := migrator.Plan{
		Migrations: []migrator.Migration{{OldKey: "GHOST", NewKey: "PHANTOM"}},
	}
	_, result, _ := migrator.Apply(s, plan)
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestApply_ConflictsOnExistingNewKey(t *testing.T) {
	s := makeSchema(ev("A", "a", true), ev("B", "b", false))
	plan := migrator.Plan{
		Migrations: []migrator.Migration{{OldKey: "A", NewKey: "B"}},
	}
	_, result, _ := migrator.Apply(s, plan)
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(result.Conflicts))
	}
}

func TestApply_UpdatesDescriptionWhenProvided(t *testing.T) {
	s := makeSchema(ev("FOO", "original", false))
	plan := migrator.Plan{
		Migrations: []migrator.Migration{{OldKey: "FOO", NewKey: "BAR", Description: "updated"}},
	}
	out, _, _ := migrator.Apply(s, plan)
	if out.Vars[0].Description != "updated" {
		t.Errorf("expected description 'updated', got %q", out.Vars[0].Description)
	}
}

func TestApply_InvalidNewKeyReturnsConflict(t *testing.T) {
	s := makeSchema(ev("VALID", "v", true))
	plan := migrator.Plan{
		Migrations: []migrator.Migration{{OldKey: "VALID", NewKey: "INVALID KEY!"}},
	}
	_, result, _ := migrator.Apply(s, plan)
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict for invalid key, got %d", len(result.Conflicts))
	}
}

func TestApply_EmptyPlanIsNoop(t *testing.T) {
	s := makeSchema(ev("X", "x", false))
	out, result, err := migrator.Apply(s, migrator.Plan{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 0 || len(result.Skipped) != 0 {
		t.Errorf("expected no-op result")
	}
	if len(out.Vars) != 1 {
		t.Errorf("expected schema unchanged")
	}
}
