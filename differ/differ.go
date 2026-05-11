package differ

import (
	"fmt"

	"github.com/yourorg/envlock/schema"
)

// Diff represents the difference between two schema versions.
type Diff struct {
	Added   []schema.EnvVar
	Removed []schema.EnvVar
	Changed []Change
}

// Change captures a modification to an existing env var definition.
type Change struct {
	Key  string
	From schema.EnvVar
	To   schema.EnvVar
}

// Compare returns a Diff between a base schema and a new schema.
func Compare(base, next *schema.Schema) *Diff {
	diff := &Diff{}

	baseMap := indexByKey(base.Vars)
	nextMap := indexByKey(next.Vars)

	for key, nextVar := range nextMap {
		baseVar, exists := baseMap[key]
		if !exists {
			diff.Added = append(diff.Added, nextVar)
			continue
		}
		if changed(baseVar, nextVar) {
			diff.Changed = append(diff.Changed, Change{
				Key:  key,
				From: baseVar,
				To:   nextVar,
			})
		}
	}

	for key, baseVar := range baseMap {
		if _, exists := nextMap[key]; !exists {
			diff.Removed = append(diff.Removed, baseVar)
		}
	}

	return diff
}

// HasChanges returns true if the diff contains any differences.
func (d *Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d *Diff) Summary() string {
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed",
		len(d.Added), len(d.Removed), len(d.Changed))
}

func indexByKey(vars []schema.EnvVar) map[string]schema.EnvVar {
	m := make(map[string]schema.EnvVar, len(vars))
	for _, v := range vars {
		m[v.Key] = v
	}
	return m
}

func changed(a, b schema.EnvVar) bool {
	return a.Required != b.Required ||
		a.Default != b.Default ||
		a.Pattern != b.Pattern ||
		a.Description != b.Description
}
