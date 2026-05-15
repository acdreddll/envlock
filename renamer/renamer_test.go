package renamer_test

import (
	"testing"

	"github.com/user/envlock/renamer"
	"github.com/user/envlock/schema"
)

func makeSchema(keys ...string) schema.Schema {
	vars := make([]schema.EnvVar, len(keys))
	for i, k := range keys {
		vars[i] = schema.EnvVar{
			Key:         k,
			Description: "desc for " + k,
			Required:    true,
		}
	}
	return schema.Schema{Vars: vars}
}

func TestRename_Success(t *testing.T) {
	s := makeSchema("OLD_KEY", "OTHER_KEY")
	newSchema, result, err := renamer.Rename(s, "OLD_KEY", "NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Applied {
		t.Errorf("expected Applied=true, got false")
	}
	if newSchema.Vars[0].Key != "NEW_KEY" {
		t.Errorf("expected first key to be NEW_KEY, got %q", newSchema.Vars[0].Key)
	}
	// Original schema must not be mutated.
	if s.Vars[0].Key != "OLD_KEY" {
		t.Errorf("original schema mutated: expected OLD_KEY, got %q", s.Vars[0].Key)
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	s := makeSchema("ALPHA", "BETA")
	_, result, err := renamer.Rename(s, "MISSING", "GAMMA")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if result.Applied {
		t.Errorf("expected Applied=false")
	}
}

func TestRename_NewKeyAlreadyExists(t *testing.T) {
	s := makeSchema("FOO", "BAR")
	_, _, err := renamer.Rename(s, "FOO", "BAR")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRename_InvalidKeyName(t *testing.T) {
	s := makeSchema("VALID_KEY")
	_, _, err := renamer.Rename(s, "VALID_KEY", "invalid-key")
	if err == nil {
		t.Fatal("expected error for invalid key name")
	}
}

func TestRename_PreservesOtherVars(t *testing.T) {
	s := makeSchema("A", "B", "C")
	newSchema, _, err := renamer.Rename(s, "B", "B_RENAMED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := make([]string, len(newSchema.Vars))
	for i, v := range newSchema.Vars {
		keys[i] = v.Key
	}
	expected := []string{"A", "B_RENAMED", "C"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}
