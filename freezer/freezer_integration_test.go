package freezer_test

import (
	"strings"
	"testing"

	"github.com/user/envlock/freezer"
	"github.com/user/envlock/schema"
)

func TestFreeze_FullWorkflow(t *testing.T) {
	s := makeSchema(
		ev("APP_NAME", required),
		ev("APP_ENV", withDefault("production")),
		ev("DB_PASS", required, sensitive),
		ev("LOG_LEVEL", withDefault("info")),
		ev("OPTIONAL_FEATURE"),
	)

	env := map[string]string{
		"APP_NAME": "envlock",
		"DB_PASS":  "s3cr3t",
	}

	snap, err := freezer.Freeze(s, env, "envlock.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// OPTIONAL_FEATURE has no value and no default — should be omitted
	if len(snap.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(snap.Entries))
	}

	redacted := snap.Redacted()
	for _, e := range redacted.Entries {
		if e.Sensitive && e.Value != "***" {
			t.Errorf("sensitive key %q not redacted", e.Key)
		}
		if !e.Sensitive && e.Value == "***" {
			t.Errorf("non-sensitive key %q was redacted", e.Key)
		}
	}

	// Validate schema file is captured
	if !strings.Contains(snap.SchemaFile, "envlock") {
		t.Errorf("unexpected schema file: %s", snap.SchemaFile)
	}

	// Defaults should be marked
	for _, e := range snap.Entries {
		switch e.Key {
		case "APP_ENV", "LOG_LEVEL":
			if !e.FromDefault {
				t.Errorf("key %q should be marked as from default", e.Key)
			}
		case "APP_NAME", "DB_PASS":
			if e.FromDefault {
				t.Errorf("key %q should not be marked as from default", e.Key)
			}
		}
	}
}
