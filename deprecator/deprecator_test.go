package deprecator_test

import (
	"testing"
	"time"

	"github.com/envlock/deprecator"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, deprecated, removeBy string) schema.EnvVar {
	return schema.EnvVar{Key: key, Deprecated: deprecated, RemoveBy: removeBy}
}

func TestDeprecate_NoDeprecations(t *testing.T) {
	s := makeSchema(
		schema.EnvVar{Key: "PORT"},
		schema.EnvVar{Key: "HOST"},
	)
	res := deprecator.Deprecate(s)
	if res.HasIssues() {
		t.Fatalf("expected no findings, got %d", len(res.Findings))
	}
}

func TestDeprecate_SingleDeprecated(t *testing.T) {
	s := makeSchema(ev("OLD_API_KEY", "2023-01-01", ""))
	res := deprecator.Deprecate(s)
	if !res.HasIssues() {
		t.Fatal("expected findings")
	}
	if res.Findings[0].Key != "OLD_API_KEY" {
		t.Errorf("expected key OLD_API_KEY, got %s", res.Findings[0].Key)
	}
}

func TestDeprecate_WithFutureRemoveBy(t *testing.T) {
	future := time.Now().AddDate(0, 6, 0).Format("2006-01-02")
	s := makeSchema(ev("LEGACY_TOKEN", "2024-01-01", future))
	res := deprecator.Deprecate(s)
	if len(res.Findings) != 1 {
		t.Fatalf("expected 1 finding")
	}
	f := res.Findings[0]
	if f.RemoveBy != future {
		t.Errorf("expected RemoveBy %s, got %s", future, f.RemoveBy)
	}
}

func TestDeprecate_PastDeadline(t *testing.T) {
	s := makeSchema(ev("DEAD_VAR", "2020-01-01", "2021-01-01"))
	res := deprecator.Deprecate(s)
	if len(res.Findings) != 1 {
		t.Fatalf("expected 1 finding")
	}
	msg := res.Findings[0].Message
	if msg == "" {
		t.Error("expected non-empty message")
	}
	// Should mention PAST deadline
	if len(msg) < 10 {
		t.Errorf("message too short: %q", msg)
	}
}

func TestDeprecate_MixedSchema(t *testing.T) {
	s := makeSchema(
		schema.EnvVar{Key: "ACTIVE"},
		ev("OLD_ONE", "2022-06-01", ""),
		ev("OLD_TWO", "2021-03-15", "2022-01-01"),
	)
	res := deprecator.Deprecate(s)
	if len(res.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(res.Findings))
	}
}
