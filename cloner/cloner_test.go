package cloner_test

import (
	"testing"

	"github.com/envlock/cloner"
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

func TestClone_Success(t *testing.T) {
	s := makeSchema(ev("DATABASE_URL", "Primary DB", "", true, true))
	out, res := cloner.Clone(s, "DATABASE_URL", "REPLICA_URL", "Replica DB", "")
	if !res.Cloned {
		t.Fatalf("expected cloned=true, got reason: %s", res.Reason)
	}
	if len(out.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(out.Vars))
	}
	cloned := out.Vars[1]
	if cloned.Key != "REPLICA_URL" {
		t.Errorf("expected key REPLICA_URL, got %s", cloned.Key)
	}
	if cloned.Description != "Replica DB" {
		t.Errorf("expected description override, got %s", cloned.Description)
	}
	if !cloned.Required || !cloned.Sensitive {
		t.Errorf("expected required and sensitive flags copied from source")
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	s := makeSchema(ev("APP_PORT", "Port", "8080", false, false))
	_, res := cloner.Clone(s, "MISSING_KEY", "NEW_KEY", "", "")
	if res.Cloned {
		t.Fatal("expected cloned=false for missing source")
	}
	if res.Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestClone_NewKeyAlreadyExists(t *testing.T) {
	s := makeSchema(
		ev("APP_PORT", "Port", "8080", false, false),
		ev("APP_HOST", "Host", "localhost", false, false),
	)
	_, res := cloner.Clone(s, "APP_PORT", "APP_HOST", "", "")
	if res.Cloned {
		t.Fatal("expected cloned=false when new key already exists")
	}
}

func TestClone_InvalidKeyName(t *testing.T) {
	s := makeSchema(ev("APP_PORT", "Port", "8080", false, false))
	_, res := cloner.Clone(s, "APP_PORT", "invalid-key", "", "")
	if res.Cloned {
		t.Fatal("expected cloned=false for invalid key name")
	}
}

func TestClone_DefaultOverride(t *testing.T) {
	s := makeSchema(ev("LOG_LEVEL", "Log level", "info", false, false))
	out, res := cloner.Clone(s, "LOG_LEVEL", "AUDIT_LOG_LEVEL", "", "debug")
	if !res.Cloned {
		t.Fatalf("expected cloned=true: %s", res.Reason)
	}
	if out.Vars[1].Default != "debug" {
		t.Errorf("expected default override 'debug', got %s", out.Vars[1].Default)
	}
	if out.Vars[1].Description != "Log level" {
		t.Errorf("expected original description preserved")
	}
}

func TestClone_DoesNotMutateOriginal(t *testing.T) {
	s := makeSchema(ev("APP_PORT", "Port", "8080", false, false))
	origLen := len(s.Vars)
	cloner.Clone(s, "APP_PORT", "APP_PORT_COPY", "", "")
	if len(s.Vars) != origLen {
		t.Error("original schema was mutated")
	}
}
