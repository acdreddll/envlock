package differ_test

import (
	"testing"

	"github.com/yourorg/envlock/differ"
	"github.com/yourorg/envlock/schema"
)

func makeSchema(vars []schema.EnvVar) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestCompare_Added(t *testing.T) {
	base := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: true},
	})
	next := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: true},
		{Key: "BAR", Required: false},
	})

	d := differ.Compare(base, next)
	if len(d.Added) != 1 || d.Added[0].Key != "BAR" {
		t.Errorf("expected BAR to be added, got %+v", d.Added)
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected removals or changes")
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: true},
		{Key: "BAR", Required: false},
	})
	next := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: true},
	})

	d := differ.Compare(base, next)
	if len(d.Removed) != 1 || d.Removed[0].Key != "BAR" {
		t.Errorf("expected BAR to be removed, got %+v", d.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	base := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: false, Default: "old"},
	})
	next := makeSchema([]schema.EnvVar{
		{Key: "FOO", Required: true, Default: "new"},
	})

	d := differ.Compare(base, next)
	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 change, got %d", len(d.Changed))
	}
	if d.Changed[0].Key != "FOO" {
		t.Errorf("expected FOO to be changed")
	}
	if d.Changed[0].From.Default != "old" || d.Changed[0].To.Default != "new" {
		t.Errorf("unexpected change values: %+v", d.Changed[0])
	}
}

func TestCompare_NoChanges(t *testing.T) {
	vars := []schema.EnvVar{{Key: "FOO", Required: true, Default: "val"}}
	d := differ.Compare(makeSchema(vars), makeSchema(vars))
	if d.HasChanges() {
		t.Errorf("expected no changes, got: %s", d.Summary())
	}
}

func TestCompare_EmptySchemas(t *testing.T) {
	d := differ.Compare(makeSchema(nil), makeSchema(nil))
	if d.HasChanges() {
		t.Errorf("expected no changes for two empty schemas, got: %s", d.Summary())
	}
}

func TestDiff_Summary(t *testing.T) {
	d := &differ.Diff{
		Added:   []schema.EnvVar{{Key: "A"}, {Key: "B"}},
		Removed: []schema.EnvVar{{Key: "C"}},
		Changed: []differ.Change{{Key: "D"}},
	}
	got := d.Summary()
	want := "+2 added, -1 removed, ~1 changed"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	d := &differ.Diff{}
	got := d.Summary()
	want := "+0 added, -0 removed, ~0 changed"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
