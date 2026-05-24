package deduplicator_test

import (
	"testing"

	"github.com/envlock/deduplicator"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string) schema.EnvVar {
	return schema.EnvVar{Key: key, Description: "desc for " + key}
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	s := makeSchema(ev("HOST"), ev("PORT"), ev("DB_URL"))
	r := deduplicator.Deduplicate(s)
	if len(r.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(r.Duplicates))
	}
	if len(r.Cleaned.Vars) != 3 {
		t.Fatalf("expected 3 vars in cleaned schema, got %d", len(r.Cleaned.Vars))
	}
}

func TestDeduplicate_SingleDuplicate(t *testing.T) {
	s := makeSchema(ev("HOST"), ev("PORT"), ev("HOST"))
	r := deduplicator.Deduplicate(s)
	if len(r.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(r.Duplicates))
	}
	if r.Duplicates[0].Key != "HOST" {
		t.Errorf("expected duplicate key HOST, got %s", r.Duplicates[0].Key)
	}
	if r.Duplicates[0].Count != 2 {
		t.Errorf("expected count 2, got %d", r.Duplicates[0].Count)
	}
	if len(r.Cleaned.Vars) != 2 {
		t.Fatalf("expected 2 vars after dedup, got %d", len(r.Cleaned.Vars))
	}
}

func TestDeduplicate_MultipleDuplicates(t *testing.T) {
	s := makeSchema(ev("A"), ev("B"), ev("A"), ev("C"), ev("B"), ev("B"))
	r := deduplicator.Deduplicate(s)
	if len(r.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(r.Duplicates))
	}
	// sorted: A, B
	if r.Duplicates[0].Key != "A" || r.Duplicates[1].Key != "B" {
		t.Errorf("unexpected duplicate order: %v", r.Duplicates)
	}
	if r.Duplicates[1].Count != 3 {
		t.Errorf("expected B count 3, got %d", r.Duplicates[1].Count)
	}
	if len(r.Cleaned.Vars) != 3 {
		t.Fatalf("expected 3 vars after dedup, got %d", len(r.Cleaned.Vars))
	}
}

func TestDeduplicate_PreservesFirstOccurrence(t *testing.T) {
	first := schema.EnvVar{Key: "TOKEN", Description: "first", Default: "abc"}
	second := schema.EnvVar{Key: "TOKEN", Description: "second", Default: "xyz"}
	s := makeSchema(first, second)
	r := deduplicator.Deduplicate(s)
	if len(r.Cleaned.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(r.Cleaned.Vars))
	}
	if r.Cleaned.Vars[0].Default != "abc" {
		t.Errorf("expected first occurrence default 'abc', got %s", r.Cleaned.Vars[0].Default)
	}
}

func TestSummary_NoDuplicates(t *testing.T) {
	r := deduplicator.Result{}
	got := deduplicator.Summary(r)
	if got != "no duplicate keys found" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithDuplicates(t *testing.T) {
	r := deduplicator.Result{
		Duplicates: []deduplicator.Duplicate{{Key: "FOO", Count: 3}},
	}
	got := deduplicator.Summary(r)
	if got == "" {
		t.Error("expected non-empty summary")
	}
}
