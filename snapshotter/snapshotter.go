package snapshotter

import (
	"fmt"
	"sort"
	"time"

	"github.com/envlock/schema"
)

// Snapshot captures the state of a schema at a point in time.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Label     string             `json:"label,omitempty"`
	Vars      []schema.EnvVar    `json:"vars"`
}

// DriftEntry describes a change between two snapshots.
type DriftEntry struct {
	Key    string `json:"key"`
	Change string `json:"change"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}

// Take creates a Snapshot from the given schema, optionally labelled.
func Take(s schema.Schema, label string) Snapshot {
	vars := make([]schema.EnvVar, len(s.Vars))
	copy(vars, s.Vars)
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Key < vars[j].Key
	})
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Vars:      vars,
	}
}

// Diff compares two snapshots and returns a list of drift entries.
func Diff(before, after Snapshot) []DriftEntry {
	var entries []DriftEntry

	before_idx := indexByKey(before.Vars)
	after_idx := indexByKey(after.Vars)

	for key, bv := range before_idx {
		av, ok := after_idx[key]
		if !ok {
			entries = append(entries, DriftEntry{Key: key, Change: "removed"})
			continue
		}
		if diffs := fieldDrift(bv, av); len(diffs) > 0 {
			entries = append(entries, diffs...)
		}
	}

	for key := range after_idx {
		if _, ok := before_idx[key]; !ok {
			entries = append(entries, DriftEntry{Key: key, Change: "added"})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Key == entries[j].Key {
			return entries[i].Change < entries[j].Change
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

func indexByKey(vars []schema.EnvVar) map[string]schema.EnvVar {
	m := make(map[string]schema.EnvVar, len(vars))
	for _, v := range vars {
		m[v.Key] = v
	}
	return m
}

func fieldDrift(b, a schema.EnvVar) []DriftEntry {
	var out []DriftEntry
	if b.Description != a.Description {
		out = append(out, DriftEntry{Key: a.Key, Change: "description", From: b.Description, To: a.Description})
	}
	if b.Default != a.Default {
		out = append(out, DriftEntry{Key: a.Key, Change: "default", From: fmt.Sprintf("%v", b.Default), To: fmt.Sprintf("%v", a.Default)})
	}
	if b.Required != a.Required {
		out = append(out, DriftEntry{Key: a.Key, Change: "required", From: fmt.Sprintf("%v", b.Required), To: fmt.Sprintf("%v", a.Required)})
	}
	return out
}
