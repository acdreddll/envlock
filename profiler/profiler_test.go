package profiler_test

import (
	"testing"

	"github.com/envlock/profiler"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc, def, pattern, group string, required, sensitive bool, tags []string) schema.EnvVar {
	return schema.EnvVar{
		Key: key, Description: desc, Default: def,
		Pattern: pattern, Group: group,
		Required: required, Sensitive: sensitive, Tags: tags,
	}
}

func TestProfile_EmptySchema(t *testing.T) {
	r := profiler.Profile(makeSchema())
	if r.TotalFields != 0 {
		t.Fatalf("expected 0 fields, got %d", r.TotalFields)
	}
	if r.AvgCompleteness != 0 {
		t.Fatalf("expected 0 avg, got %f", r.AvgCompleteness)
	}
}

func TestProfile_FullyAnnotated(t *testing.T) {
	s := makeSchema(ev("PORT", "HTTP port", "8080", `^\d+$`, "network", false, false, []string{"infra"}))
	r := profiler.Profile(s)
	if r.TotalFields != 1 {
		t.Fatalf("expected 1 field")
	}
	if r.Fields[0].CompletenessScore != 100 {
		t.Errorf("expected score 100, got %d", r.Fields[0].CompletenessScore)
	}
	if len(r.MissingDesc) != 0 {
		t.Errorf("expected no missing desc")
	}
}

func TestProfile_MissingDescription(t *testing.T) {
	s := makeSchema(ev("SECRET", "", "", "", "", true, true, nil))
	r := profiler.Profile(s)
	if len(r.MissingDesc) != 1 || r.MissingDesc[0] != "SECRET" {
		t.Errorf("expected SECRET in MissingDesc, got %v", r.MissingDesc)
	}
}

func TestProfile_MissingDefaultOptional(t *testing.T) {
	s := makeSchema(ev("TIMEOUT", "timeout", "", "", "", false, false, nil))
	r := profiler.Profile(s)
	if len(r.MissingDefault) != 1 || r.MissingDefault[0] != "TIMEOUT" {
		t.Errorf("expected TIMEOUT in MissingDefault, got %v", r.MissingDefault)
	}
}

func TestProfile_RequiredSkipsMissingDefault(t *testing.T) {
	s := makeSchema(ev("DB_URL", "db", "", "", "", true, false, nil))
	r := profiler.Profile(s)
	if len(r.MissingDefault) != 0 {
		t.Errorf("required field should not appear in MissingDefault")
	}
}

func TestProfile_AvgCompleteness(t *testing.T) {
	s := makeSchema(
		ev("A", "desc", "val", `.*`, "g", false, false, []string{"t"}), // 100
		ev("B", "", "", "", "", false, false, nil),                       // 0
	)
	r := profiler.Profile(s)
	if r.AvgCompleteness != 50.0 {
		t.Errorf("expected avg 50.0, got %f", r.AvgCompleteness)
	}
}

func TestSummary_Format(t *testing.T) {
	s := makeSchema(ev("X", "x", "1", "", "", false, false, nil))
	r := profiler.Profile(s)
	got := profiler.Summary(r)
	if got == "" {
		t.Error("expected non-empty summary")
	}
}
