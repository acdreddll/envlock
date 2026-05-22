package pruner_test

import (
	"testing"

	"github.com/your-org/envlock/pruner"
	"github.com/your-org/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, required bool, def, desc, deprecated string) schema.EnvVar {
	return schema.EnvVar{
		Key:         key,
		Required:    required,
		Default:     def,
		Description: desc,
		Deprecated:  deprecated,
	}
}

func TestPrune_NoOptions(t *testing.T) {
	s := makeSchema(ev("A", true, "", "desc", ""), ev("B", false, "", "", ""))
	res := pruner.Prune(s, pruner.Options{})
	if len(res.Kept) != 2 || len(res.Removed) != 0 {
		t.Fatalf("expected 2 kept, 0 removed; got %d kept, %d removed", len(res.Kept), len(res.Removed))
	}
}

func TestPrune_RemoveDeprecated(t *testing.T) {
	s := makeSchema(
		ev("ACTIVE", true, "", "still used", ""),
		ev("OLD_KEY", false, "", "old", "use NEW_KEY instead"),
	)
	res := pruner.Prune(s, pruner.Options{RemoveDeprecated: true})
	if len(res.Kept) != 1 || res.Kept[0].Key != "ACTIVE" {
		t.Fatalf("expected ACTIVE kept; got %+v", res.Kept)
	}
	if len(res.Removed) != 1 || res.Removed[0].Key != "OLD_KEY" {
		t.Fatalf("expected OLD_KEY removed; got %+v", res.Removed)
	}
}

func TestPrune_RemoveOptionalNoDefault(t *testing.T) {
	s := makeSchema(
		ev("REQUIRED", true, "", "", ""),
		ev("HAS_DEFAULT", false, "val", "", ""),
		ev("HAS_DESC", false, "", "something", ""),
		ev("BARE", false, "", "", ""),
	)
	res := pruner.Prune(s, pruner.Options{RemoveOptionalNoDefault: true})
	if len(res.Removed) != 1 || res.Removed[0].Key != "BARE" {
		t.Fatalf("expected only BARE removed; got %+v", res.Removed)
	}
}

func TestPrune_ExplicitKeys(t *testing.T) {
	s := makeSchema(
		ev("A", true, "", "a", ""),
		ev("B", true, "", "b", ""),
		ev("C", false, "", "c", ""),
	)
	res := pruner.Prune(s, pruner.Options{Keys: []string{"A", "C"}})
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed; got %d", len(res.Removed))
	}
	if len(res.Kept) != 1 || res.Kept[0].Key != "B" {
		t.Fatalf("expected only B kept; got %+v", res.Kept)
	}
}

func TestPrune_ExplicitKeysTakePrecedence(t *testing.T) {
	// Even if RemoveDeprecated is false, explicit key list still removes.
	s := makeSchema(ev("LEGACY", false, "", "", "deprecated"))
	res := pruner.Prune(s, pruner.Options{RemoveDeprecated: false, Keys: []string{"LEGACY"}})
	if len(res.Removed) != 1 {
		t.Fatalf("expected LEGACY removed via explicit key; got %+v", res.Removed)
	}
}
