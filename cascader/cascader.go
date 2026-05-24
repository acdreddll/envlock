// Package cascader resolves environment variable values by cascading through
// multiple sources in priority order: overrides > env > defaults.
package cascader

import (
	"fmt"

	"github.com/envlock/schema"
)

// Source represents a named map of key→value pairs.
type Source struct {
	Name   string
	Values map[string]string
}

// Result holds the resolved value for a single variable.
type Result struct {
	Key      string
	Value    string
	Source   string
	Resolved bool
}

// Cascade resolves each variable in the schema by walking sources in order
// (first source wins). If no source provides a value, the schema default is
// used. Required variables with no resolved value produce an error entry.
func Cascade(s schema.Schema, sources []Source) ([]Result, []error) {
	var results []Result
	var errs []error

	for _, v := range s.Vars {
		r := Result{Key: v.Key}

		// Walk sources in priority order.
		for _, src := range sources {
			if val, ok := src.Values[v.Key]; ok && val != "" {
				r.Value = val
				r.Source = src.Name
				r.Resolved = true
				break
			}
		}

		// Fall back to schema default.
		if !r.Resolved && v.Default != "" {
			r.Value = v.Default
			r.Source = "default"
			r.Resolved = true
		}

		if !r.Resolved && v.Required {
			errs = append(errs, fmt.Errorf("required variable %q not found in any source", v.Key))
		}

		results = append(results, r)
	}

	return results, errs
}

// ToMap converts a slice of Results into a plain key→value map,
// including only resolved entries.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		if r.Resolved {
			out[r.Key] = r.Value
		}
	}
	return out
}
