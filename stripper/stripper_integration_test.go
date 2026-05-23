package stripper_test

import (
	"testing"

	"github.com/your-org/envlock/schema"
	"github.com/your-org/envlock/stripper"
)

// TestStrip_PublishWorkflow simulates stripping sensitive metadata before
// publishing a schema to an external team.
func TestStrip_PublishWorkflow(t *testing.T) {
	s := schema.Schema{
		{Key: "DB_PASS", Description: "DB password", Sensitive: true, Required: true, Pattern: `^.{8,}$`, Group: "db", Tags: []string{"secret"}},
		{Key: "APP_ENV", Description: "Environment name", Default: "production", Group: "app"},
		{Key: "OLD_KEY", Description: "Legacy key", Deprecated: true, RemoveBy: "2024-06-01"},
	}

	out, err := stripper.Strip(s, stripper.Options{
		Fields: []string{"sensitive", "pattern", "deprecated", "remove_by"},
	})
	if err != nil {
		t.Fatalf("Strip failed: %v", err)
	}

	if len(out) != 3 {
		t.Fatalf("expected 3 vars, got %d", len(out))
	}

	for _, v := range out {
		if v.Sensitive {
			t.Errorf("%s: sensitive should be stripped", v.Key)
		}
		if v.Pattern != "" {
			t.Errorf("%s: pattern should be stripped", v.Key)
		}
		if v.Deprecated {
			t.Errorf("%s: deprecated should be stripped", v.Key)
		}
		if v.RemoveBy != "" {
			t.Errorf("%s: remove_by should be stripped", v.Key)
		}
	}

	// Descriptions and defaults should be preserved
	if out[0].Description == "" {
		t.Error("DB_PASS description should be preserved")
	}
	if out[1].Default == "" {
		t.Error("APP_ENV default should be preserved")
	}
}
