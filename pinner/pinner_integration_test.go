package pinner_test

import (
	"testing"

	"github.com/envlock/pinner"
)

func TestPinnerWorkflow_FullRoundTrip(t *testing.T) {
	s := makeSchema("API_KEY", "DB_HOST", "TIMEOUT")

	originalEnv := map[string]string{
		"API_KEY": "secret-abc",
		"DB_HOST": "db.internal",
		"TIMEOUT": "30s",
	}

	snap := pinner.Capture(s, originalEnv)
	if len(snap.Pins) != 3 {
		t.Fatalf("expected 3 pins, got %d", len(snap.Pins))
	}

	// No drift when env unchanged.
	report := pinner.Detect(snap, originalEnv)
	if report.HasDrift() {
		t.Errorf("expected clean, got drift: %v", report.Drifted)
	}

	// Simulate a deployment that changes DB_HOST and removes TIMEOUT.
	changedEnv := map[string]string{
		"API_KEY": "secret-abc",
		"DB_HOST": "db.prod",
	}

	report = pinner.Detect(snap, changedEnv)
	if !report.HasDrift() {
		t.Fatal("expected drift after env change")
	}

	if len(report.Drifted) != 2 {
		t.Errorf("expected 2 drifted entries, got %d", len(report.Drifted))
	}

	driftedKeys := map[string]bool{}
	for _, d := range report.Drifted {
		driftedKeys[d.Key] = true
	}
	if !driftedKeys["DB_HOST"] {
		t.Error("expected DB_HOST to be drifted")
	}
	if !driftedKeys["TIMEOUT"] {
		t.Error("expected TIMEOUT to be missing/drifted")
	}

	if len(report.Clean) != 1 || report.Clean[0] != "API_KEY" {
		t.Errorf("expected API_KEY to be clean, got %v", report.Clean)
	}
}
