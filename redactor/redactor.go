package redactor

import (
	"strings"

	"github.com/envlock/schema"
)

const defaultMask = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// Mask is the string used to replace sensitive values.
	// Defaults to "***REDACTED***" when empty.
	Mask string
	// RedactAll replaces every value regardless of the sensitive flag.
	RedactAll bool
}

// Result holds the redacted key/value pairs.
type Result struct {
	Key      string
	Value    string
	Redacted bool
}

// Redact takes a schema and a map of resolved env values and returns a slice
// of Result entries where sensitive values have been masked.
func Redact(s schema.Schema, env map[string]string, opts Options) []Result {
	mask := opts.Mask
	if mask == "" {
		mask = defaultMask
	}

	results := make([]Result, 0, len(s.Vars))
	for _, v := range s.Vars {
		raw, ok := env[v.Key]
		if !ok {
			if v.Default != "" {
				raw = v.Default
			}
		}

		should := opts.RedactAll || v.Sensitive
		if should && strings.TrimSpace(raw) != "" {
			results = append(results, Result{Key: v.Key, Value: mask, Redacted: true})
		} else {
			results = append(results, Result{Key: v.Key, Value: raw, Redacted: false})
		}
	}
	return results
}

// IndexByKey converts a []Result slice into a map keyed by variable name for
// convenient lookup.
func IndexByKey(results []Result) map[string]Result {
	m := make(map[string]Result, len(results))
	for _, r := range results {
		m[r.Key] = r
	}
	return m
}
