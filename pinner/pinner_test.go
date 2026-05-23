package pinner_test

import (
	"testing"

	"github.com/envlock/pinner"
	"github.com/envlock/schema"
)

func makeSchema(keys ...string) schema.Schema {
	var vars []schema.EnvVar
	for _, k := range keys {
		vars = append(vars, schema.EnvVar{Key: k})
	}
	return schema.Schema{Vars: vars}
}

func TestCapture_AllPresent(t *testing.T) {
	s := makeSchema("HOST", "PORT", "DB_URL")
	env := map[string]string{"HOST": "localhost", "PORT": "5432", "DB_URL": "postgres://"}
	snap := pinner.Capture(s, env)
	if len(snap.Pins) != 3 {
		t.Fatalf("expected 3 pins, got %d", len(snap.Pins))
	}
}

func TestCapture_SkipsMissingKeys(t *testing.T) {
	s := makeSchema("HOST", "PORT")
	env := map[string]string{"HOST": "localhost"}
	snap := pinner.Capture(s, env)
	if len(snap.Pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(snap.Pins))
	}
	if snap.Pins[0].Key != "HOST" {
		t.Errorf("unexpected key %s", snap.Pins[0].Key)
	}
}

func TestDetect_NoDrift(t *testing.T) {
	s := makeSchema("HOST", "PORT")
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	snap := pinner.Capture(s, env)
	report := pinner.Detect(snap, env)
	if report.HasDrift() {
		t.Errorf("expected no drift, got %+v", report.Drifted)
	}
	if len(report.Clean) != 2 {
		t.Errorf("expected 2 clean vars, got %d", len(report.Clean))
	}
}

func TestDetect_ValueChanged(t *testing.T) {
	s := makeSchema("HOST")
	snap := pinner.Capture(s, map[string]string{"HOST": "localhost"})
	report := pinner.Detect(snap, map[string]string{"HOST": "remotehost"})
	if !report.HasDrift() {
		t.Fatal("expected drift")
	}
	if report.Drifted[0].Pinned != "localhost" || report.Drifted[0].Current != "remotehost" {
		t.Errorf("unexpected drift entry: %+v", report.Drifted[0])
	}
}

func TestDetect_KeyRemoved(t *testing.T) {
	s := makeSchema("HOST")
	snap := pinner.Capture(s, map[string]string{"HOST": "localhost"})
	report := pinner.Detect(snap, map[string]string{})
	if !report.HasDrift() {
		t.Fatal("expected drift due to missing key")
	}
	if !report.Drifted[0].Missing {
		t.Errorf("expected Missing=true, got %+v", report.Drifted[0])
	}
}

func TestSummary_NoDrift(t *testing.T) {
	r := pinner.DriftReport{Clean: []string{"A", "B"}}
	s := pinner.Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_WithDrift(t *testing.T) {
	r := pinner.DriftReport{
		Drifted: []pinner.DriftEntry{{Key: "X", Pinned: "old", Current: "new"}},
		Clean:   []string{"Y"},
	}
	s := pinner.Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
