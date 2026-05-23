// Package pinner locks environment variable values to a snapshot, detecting
// drift between the pinned snapshot and a live environment.
package pinner

import (
	"fmt"
	"sort"
	"time"

	"github.com/envlock/schema"
)

// Pin represents a single pinned variable value.
type Pin struct {
	Key       string    `json:"key" yaml:"key"`
	Value     string    `json:"value" yaml:"value"`
	PinnedAt  time.Time `json:"pinned_at" yaml:"pinned_at"`
}

// Snapshot is a collection of pinned variables.
type Snapshot struct {
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	Pins      []Pin     `json:"pins" yaml:"pins"`
}

// DriftEntry describes a variable that has changed since pinning.
type DriftEntry struct {
	Key      string
	Pinned   string
	Current  string
	Missing  bool
}

// DriftReport holds the result of comparing a snapshot to live env values.
type DriftReport struct {
	Drifted []DriftEntry
	Clean   []string
}

// HasDrift returns true when at least one variable has drifted.
func (r DriftReport) HasDrift() bool { return len(r.Drifted) > 0 }

// Capture builds a Snapshot from the provided env map, scoped to keys present
// in the schema.
func Capture(s schema.Schema, env map[string]string) Snapshot {
	now := time.Now().UTC()
	var pins []Pin
	for _, v := range s.Vars {
		val, ok := env[v.Key]
		if !ok {
			continue
		}
		pins = append(pins, Pin{Key: v.Key, Value: val, PinnedAt: now})
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Key < pins[j].Key })
	return Snapshot{CreatedAt: now, Pins: pins}
}

// Detect compares a previously captured Snapshot to a live env map and returns
// a DriftReport describing any changes.
func Detect(snap Snapshot, env map[string]string) DriftReport {
	var report DriftReport
	for _, p := range snap.Pins {
		current, ok := env[p.Key]
		if !ok {
			report.Drifted = append(report.Drifted, DriftEntry{
				Key:     p.Key,
				Pinned:  p.Value,
				Missing: true,
			})
			continue
		}
		if current != p.Value {
			report.Drifted = append(report.Drifted, DriftEntry{
				Key:     p.Key,
				Pinned:  p.Value,
				Current: current,
			})
		} else {
			report.Clean = append(report.Clean, p.Key)
		}
	}
	return report
}

// Summary returns a human-readable one-line summary of the drift report.
func Summary(r DriftReport) string {
	if !r.HasDrift() {
		return fmt.Sprintf("no drift detected (%d pinned vars clean)", len(r.Clean))
	}
	return fmt.Sprintf("%d var(s) drifted, %d clean", len(r.Drifted), len(r.Clean))
}
