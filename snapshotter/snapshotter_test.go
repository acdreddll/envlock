package snapshotter_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/snapshotter"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc, def string, required bool) schema.EnvVar {
	return schema.EnvVar{Key: key, Description: desc, Default: def, Required: required}
}

func TestTake_SortsVars(t *testing.T) {
	s := makeSchema(
		ev("Z_VAR", "last", "", false),
		ev("A_VAR", "first", "", true),
	)
	snap := snapshotter.Take(s, "test")
	if snap.Vars[0].Key != "A_VAR" {
		t.Errorf("expected A_VAR first, got %s", snap.Vars[0].Key)
	}
	if snap.Label != "test" {
		t.Errorf("expected label 'test', got %s", snap.Label)
	}
}

func TestTake_TimestampSet(t *testing.T) {
	snap := snapshotter.Take(makeSchema(ev("X", "", "", false)), "")
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestDiff_NoDrift(t *testing.T) {
	s := makeSchema(ev("FOO", "desc", "val", true))
	before := snapshotter.Take(s, "before")
	after := snapshotter.Take(s, "after")
	entries := snapshotter.Diff(before, after)
	if len(entries) != 0 {
		t.Errorf("expected no drift, got %d entries", len(entries))
	}
}

func TestDiff_Added(t *testing.T) {
	before := snapshotter.Take(makeSchema(ev("A", "", "", false)), "")
	after := snapshotter.Take(makeSchema(ev("A", "", "", false), ev("B", "", "", false)), "")
	entries := snapshotter.Diff(before, after)
	if len(entries) != 1 || entries[0].Key != "B" || entries[0].Change != "added" {
		t.Errorf("unexpected drift entries: %+v", entries)
	}
}

func TestDiff_Removed(t *testing.T) {
	before := snapshotter.Take(makeSchema(ev("A", "", "", false), ev("B", "", "", false)), "")
	after := snapshotter.Take(makeSchema(ev("A", "", "", false)), "")
	entries := snapshotter.Diff(before, after)
	if len(entries) != 1 || entries[0].Key != "B" || entries[0].Change != "removed" {
		t.Errorf("unexpected drift entries: %+v", entries)
	}
}

func TestDiff_FieldChanged(t *testing.T) {
	before := snapshotter.Take(makeSchema(ev("FOO", "old desc", "old", false)), "")
	after := snapshotter.Take(makeSchema(ev("FOO", "new desc", "new", false)), "")
	entries := snapshotter.Diff(before, after)
	if len(entries) != 2 {
		t.Fatalf("expected 2 field diffs, got %d: %+v", len(entries), entries)
	}
	changes := map[string]bool{}
	for _, e := range entries {
		changes[e.Change] = true
	}
	if !changes["description"] || !changes["default"] {
		t.Errorf("expected description and default changes, got %+v", entries)
	}
}

func TestDiff_RequiredChanged(t *testing.T) {
	before := snapshotter.Take(makeSchema(ev("BAR", "", "", false)), "")
	after := snapshotter.Take(makeSchema(ev("BAR", "", "", true)), "")
	entries := snapshotter.Diff(before, after)
	if len(entries) != 1 || entries[0].Change != "required" {
		t.Errorf("expected required change, got %+v", entries)
	}
}
