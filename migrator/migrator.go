package migrator

import (
	"fmt"
	"strings"

	"github.com/envlock/schema"
)

// Migration represents a single rename/transform operation.
type Migration struct {
	OldKey      string `yaml:"old_key"`
	NewKey      string `yaml:"new_key"`
	Description string `yaml:"description,omitempty"`
}

// Plan holds all migrations to apply.
type Plan struct {
	Migrations []Migration `yaml:"migrations"`
}

// Result describes the outcome of applying a migration plan.
type Result struct {
	Applied  []Migration
	Skipped  []Migration
	Conflicts []string
}

// Apply executes the migration plan against the given schema, returning a new
// schema with keys renamed according to the plan.
func Apply(s schema.Schema, plan Plan) (schema.Schema, Result, error) {
	result := Result{}

	// Build a mutable copy keyed by name.
	index := make(map[string]schema.EnvVar, len(s.Vars))
	for _, ev := range s.Vars {
		index[ev.Key] = ev
	}

	for _, m := range plan.Migrations {
		if err := validateKey(m.NewKey); err != nil {
			result.Conflicts = append(result.Conflicts,
				fmt.Sprintf("invalid new key %q: %v", m.NewKey, err))
			result.Skipped = append(result.Skipped, m)
			continue
		}

		ev, exists := index[m.OldKey]
		if !exists {
			result.Skipped = append(result.Skipped, m)
			continue
		}

		if _, conflict := index[m.NewKey]; conflict {
			result.Conflicts = append(result.Conflicts,
				fmt.Sprintf("key %q already exists; cannot rename %q to it", m.NewKey, m.OldKey))
			result.Skipped = append(result.Skipped, m)
			continue
		}

		delete(index, m.OldKey)
		ev.Key = m.NewKey
		if m.Description != "" {
			ev.Description = m.Description
		}
		index[m.NewKey] = ev
		result.Applied = append(result.Applied, m)
	}

	out := schema.Schema{}
	for _, ev := range index {
		out.Vars = append(out.Vars, ev)
	}
	return out, result, nil
}

func validateKey(k string) error {
	if strings.TrimSpace(k) == "" {
		return fmt.Errorf("key must not be empty")
	}
	for _, r := range k {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') || r == '_') {
			return fmt.Errorf("key contains invalid character %q", r)
		}
	}
	return nil
}
