package flattener

import (
	"fmt"
	"sort"

	"github.com/envlock/schema"
)

// FlatEntry represents a single environment variable entry in a flat key=value style.
type FlatEntry struct {
	Key         string
	Value       string
	Required    bool
	Sensitive   bool
	Description string
}

// FlatResult holds the result of a Flatten operation.
type FlatResult struct {
	Entries []FlatEntry
}

// Flatten converts a schema and a resolved env map into a flat list of FlatEntry values.
// Keys not present in env fall back to their schema default.
// If a required key is missing from both env and defaults, an error is returned.
func Flatten(s schema.Schema, env map[string]string) (*FlatResult, error) {
	entries := make([]FlatEntry, 0, len(s.Vars))

	for _, v := range s.Vars {
		val, ok := env[v.Key]
		if !ok {
			if v.Default != "" {
				val = v.Default
			} else if v.Required {
				return nil, fmt.Errorf("required key %q is missing and has no default", v.Key)
			} else {
				// optional, no value — include with empty string
				val = ""
			}
		}

		entries = append(entries, FlatEntry{
			Key:         v.Key,
			Value:       val,
			Required:    v.Required,
			Sensitive:   v.Sensitive,
			Description: v.Description,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &FlatResult{Entries: entries}, nil
}

// ToMap converts a FlatResult into a plain map[string]string.
func (r *FlatResult) ToMap() map[string]string {
	out := make(map[string]string, len(r.Entries))
	for _, e := range r.Entries {
		out[e.Key] = e.Value
	}
	return out
}
