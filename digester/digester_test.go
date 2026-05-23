package digester_test

import (
	"testing"

	"github.com/envlock/digester"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc, def string, required, sensitive bool) schema.EnvVar {
	return schema.EnvVar{
		Key:         key,
		Description: desc,
		Default:     def,
		Required:    required,
		Sensitive:   sensitive,
	}
}

func TestCompute_EmptySchema(t *testing.T) {
	d := digester.Compute(makeSchema())
	if d.KeyCount != 0 {
		t.Errorf("expected 0 keys, got %d", d.KeyCount)
	}
	if d.SchemaHash == "" {
		t.Error("expected non-empty schema hash")
	}
	if len(d.KeyHashes) != 0 {
		t.Errorf("expected empty key hashes, got %d", len(d.KeyHashes))
	}
}

func TestCompute_KeyCount(t *testing.T) {
	s := makeSchema(
		ev("A", "desc a", "", true, false),
		ev("B", "desc b", "val", false, true),
	)
	d := digester.Compute(s)
	if d.KeyCount != 2 {
		t.Errorf("expected 2, got %d", d.KeyCount)
	}
	if _, ok := d.KeyHashes["A"]; !ok {
		t.Error("expected hash for key A")
	}
	if _, ok := d.KeyHashes["B"]; !ok {
		t.Error("expected hash for key B")
	}
}

func TestCompute_Deterministic(t *testing.T) {
	s := makeSchema(
		ev("X", "some desc", "default", true, false),
		ev("Y", "other", "", false, true),
	)
	d1 := digester.Compute(s)
	d2 := digester.Compute(s)
	if d1.SchemaHash != d2.SchemaHash {
		t.Errorf("expected identical hashes, got %s vs %s", d1.SchemaHash, d2.SchemaHash)
	}
}

func TestCompute_ChangeSensitiveAffectsHash(t *testing.T) {
	s1 := makeSchema(ev("SECRET", "a secret", "", true, false))
	s2 := makeSchema(ev("SECRET", "a secret", "", true, true))
	d1 := digester.Compute(s1)
	d2 := digester.Compute(s2)
	if d1.SchemaHash == d2.SchemaHash {
		t.Error("expected different hashes when sensitive flag changes")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	s := makeSchema(ev("A", "desc", "", true, false))
	d := digester.Compute(s)
	added, removed, changed := digester.Diff(d, d)
	if len(added)+len(removed)+len(changed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v", added, removed, changed)
	}
}

func TestDiff_Added(t *testing.T) {
	a := makeSchema(ev("A", "desc", "", true, false))
	b := makeSchema(ev("A", "desc", "", true, false), ev("B", "new", "", false, false))
	da := digester.Compute(a)
	db := digester.Compute(b)
	added, removed, changed := digester.Diff(da, db)
	if len(added) != 1 || added[0] != "B" {
		t.Errorf("expected [B] added, got %v", added)
	}
	if len(removed) != 0 || len(changed) != 0 {
		t.Errorf("unexpected removed/changed: %v %v", removed, changed)
	}
}

func TestDiff_Removed(t *testing.T) {
	a := makeSchema(ev("A", "desc", "", true, false), ev("B", "old", "", false, false))
	b := makeSchema(ev("A", "desc", "", true, false))
	da := digester.Compute(a)
	db := digester.Compute(b)
	_, removed, _ := digester.Diff(da, db)
	if len(removed) != 1 || removed[0] != "B" {
		t.Errorf("expected [B] removed, got %v", removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	a := makeSchema(ev("A", "original desc", "", true, false))
	b := makeSchema(ev("A", "updated desc", "", true, false))
	da := digester.Compute(a)
	db := digester.Compute(b)
	_, _, changed := digester.Diff(da, db)
	if len(changed) != 1 || changed[0] != "A" {
		t.Errorf("expected [A] changed, got %v", changed)
	}
}
