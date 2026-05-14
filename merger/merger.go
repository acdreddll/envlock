package merger

import (
	"fmt"

	"github.com/envlock/schema"
)

// ConflictError describes a key conflict between two schemas.
type ConflictError struct {
	Key    string
	Source string
	Target string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("conflict on key %q: source required=%v, target required=%v", e.Key, e.Source, e.Target)
}

// Result holds the merged schema and any non-fatal warnings.
type Result struct {
	Schema   schema.Schema
	Warnings []string
}

// Merge combines two schemas into one. Keys present in both are merged with
// target values taking precedence. Conflicts in required/sensitive flags are
// recorded as warnings rather than hard errors.
func Merge(base, override schema.Schema) (Result, error) {
	result := Result{}

	index := make(map[string]schema.EnvVar)
	for _, ev := range base.Vars {
		index[ev.Key] = ev
	}

	for _, ov := range override.Vars {
		if existing, ok := index[ov.Key]; ok {
			merged := existing

			// Override wins on description and pattern if set.
			if ov.Description != "" {
				merged.Description = ov.Description
			}
			if ov.Pattern != "" {
				merged.Pattern = ov.Pattern
			}
			if ov.Default != "" {
				merged.Default = ov.Default
			}

			// Warn when required flag differs.
			if existing.Required != ov.Required {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("key %q: required flag differs (base=%v, override=%v); using override",
						ov.Key, existing.Required, ov.Required))
				merged.Required = ov.Required
			}

			// Warn when sensitive flag differs.
			if existing.Sensitive != ov.Sensitive {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("key %q: sensitive flag differs (base=%v, override=%v); using override",
						ov.Key, existing.Sensitive, ov.Sensitive))
				merged.Sensitive = ov.Sensitive
			}

			index[ov.Key] = merged
		} else {
			index[ov.Key] = ov
		}
	}

	// Preserve base order, then append new keys from override.
	seen := make(map[string]bool)
	for _, ev := range base.Vars {
		result.Schema.Vars = append(result.Schema.Vars, index[ev.Key])
		seen[ev.Key] = true
	}
	for _, ov := range override.Vars {
		if !seen[ov.Key] {
			result.Schema.Vars = append(result.Schema.Vars, index[ov.Key])
		}
	}

	return result, nil
}
