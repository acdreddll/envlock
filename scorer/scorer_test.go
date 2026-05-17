package scorer_test

import (
	"testing"

	"github.com/user/envlock/schema"
	"github.com/user/envlock/scorer"
)

func makeSchema(vars []schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func TestScore_EmptySchema(t *testing.T) {
	r := scorer.Score(makeSchema(nil))
	if r.Total != 0 || r.Score != 0 {
		t.Errorf("expected zero score for empty schema, got %+v", r)
	}
}

func TestScore_FullyAnnotated(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "DB_URL", Description: "Database URL", Group: "db", Required: false, Default: "localhost", Sensitive: false},
	}
	r := scorer.Score(makeSchema(vars))
	// description(40) + default(20) + group(20) = 80
	if r.Score != 80 {
		t.Errorf("expected score 80, got %d", r.Score)
	}
}

func TestScore_SensitiveWithPattern(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "API_SECRET", Description: "Secret key", Sensitive: true, Pattern: `^[A-Za-z0-9]{32}$`, Group: "auth", Required: true},
	}
	r := scorer.Score(makeSchema(vars))
	// description(40) + pattern(20) + group(20) = 80
	if r.Score != 80 {
		t.Errorf("expected score 80, got %d", r.Score)
	}
}

func TestScore_NoAnnotations(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "FOO"},
		{Key: "BAR"},
	}
	r := scorer.Score(makeSchema(vars))
	if r.Score != 0 {
		t.Errorf("expected score 0, got %d", r.Score)
	}
	if r.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Total)
	}
}

func TestScore_BreakdownCounts(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "A", Description: "desc a", Group: "g1"},
		{Key: "B", Description: "desc b"},
		{Key: "C"},
	}
	r := scorer.Score(makeSchema(vars))
	if r.Breakdown["description"] != 2 {
		t.Errorf("expected 2 descriptions, got %d", r.Breakdown["description"])
	}
	if r.Breakdown["group"] != 1 {
		t.Errorf("expected 1 group, got %d", r.Breakdown["group"])
	}
}

func TestScore_PercentageRange(t *testing.T) {
	vars := []schema.EnvVar{
		{Key: "X", Description: "x", Group: "g", Default: "v", Required: false},
	}
	r := scorer.Score(makeSchema(vars))
	if r.Percentage < 0 || r.Percentage > 100 {
		t.Errorf("percentage out of range: %f", r.Percentage)
	}
}
