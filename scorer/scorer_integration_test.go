package scorer_test

import (
	"testing"

	"github.com/user/envlock/schema"
	"github.com/user/envlock/scorer"
)

func TestScore_MixedSchema(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "DB_HOST", Description: "DB hostname", Group: "database", Required: true},
		{Key: "DB_PASS", Sensitive: true, Pattern: `^.{8,}$`, Group: "database", Required: true},
		{Key: "LOG_LEVEL", Default: "info", Required: false},
		{Key: "ORPHAN"},
	}
	r := scorer.Score(schema.Schema{Vars: vars})

	if r.Total != 4 {
		t.Fatalf("expected 4 vars, got %d", r.Total)
	}
	if r.Breakdown["description"] != 1 {
		t.Errorf("expected 1 description, got %d", r.Breakdown["description"])
	}
	if r.Breakdown["group"] != 2 {
		t.Errorf("expected 2 groups, got %d", r.Breakdown["group"])
	}
	if r.Breakdown["pattern"] != 1 {
		t.Errorf("expected 1 pattern, got %d", r.Breakdown["pattern"])
	}
	if r.Breakdown["default"] != 1 {
		t.Errorf("expected 1 default, got %d", r.Breakdown["default"])
	}
	if r.Score < 0 || r.Score > 100 {
		t.Errorf("score out of bounds: %d", r.Score)
	}
}
