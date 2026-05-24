package snapshotter_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/snapshotter"
)

func TestSnapshot_FullWorkflow(t *testing.T) {
	v1 := makeSchema(
		ev("DB_HOST", "database host", "localhost", true),
		ev("DB_PORT", "database port", "5432", false),
		ev("API_KEY", "api key", "", true),
	)

	v2 := makeSchema(
		ev("DB_HOST", "database host", "prod.db", true), // default changed
		ev("DB_PORT", "database port", "5432", false),
		ev("API_KEY", "api key", "", true),
		ev("CACHE_TTL", "cache ttl", "300", false), // added
	)

	snap1 := snapshotter.Take(v1, "v1")
	snap2 := snapshotter.Take(v2, "v2")

	if snap1.Label != "v1" {
		t.Errorf("expected label v1, got %s", snap1.Label)
	}
	if len(snap1.Vars) != 3 {
		t.Errorf("expected 3 vars in snap1, got %d", len(snap1.Vars))
	}

	drift := snapshotter.Diff(snap1, snap2)

	changeMap := map[string][]string{}
	for _, d := range drift {
		changeMap[d.Key] = append(changeMap[d.Key], d.Change)
	}

	if _, ok := changeMap["CACHE_TTL"]; !ok {
		t.Error("expected CACHE_TTL to appear as added")
	}
	if _, ok := changeMap["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to have drift")
	}
	for _, ch := range changeMap["DB_HOST"] {
		if ch == "default" {
			return
		}
	}
	t.Error("expected default change for DB_HOST")
}

func TestSnapshot_DoesNotMutateOriginal(t *testing.T) {
	orig := []schema.EnvVar{
		ev("Z", "", "", false),
		ev("A", "", "", false),
	}
	s := schema.Schema{Vars: orig}
	snapshotter.Take(s, "")
	if s.Vars[0].Key != "Z" {
		t.Error("Take mutated the original schema slice")
	}
}
