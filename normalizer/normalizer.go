package normalizer

import (
	"strings"

	"github.com/envlock/schema"
)

// Result holds the outcome of normalizing a single variable.
type Result struct {
	Key      string
	Original string
	Normalized string
	Changed  bool
}

// Options controls which normalizations are applied.
type Options struct {
	TrimSpace   bool
	UpperCase   bool
	LowerCase   bool
	CollapseWS  bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		TrimSpace:  true,
		UpperCase:  false,
		LowerCase:  false,
		CollapseWS: false,
	}
}

// Normalize applies normalization rules to the provided env map according to
// the schema and options. Only keys defined in the schema are normalized.
// Returns a slice of Result for every key that was evaluated.
func Normalize(s schema.Schema, env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(s.Vars))

	for _, v := range s.Vars {
		raw, ok := env[v.Key]
		if !ok {
			continue
		}

		normalized := raw

		if opts.TrimSpace {
			normalized = strings.TrimSpace(normalized)
		}

		if opts.CollapseWS {
			normalized = collapseWhitespace(normalized)
		}

		if opts.UpperCase {
			normalized = strings.ToUpper(normalized)
		} else if opts.LowerCase {
			normalized = strings.ToLower(normalized)
		}

		results = append(results, Result{
			Key:        v.Key,
			Original:   raw,
			Normalized: normalized,
			Changed:    raw != normalized,
		})
	}

	return results
}

// collapseWhitespace replaces runs of whitespace with a single space.
func collapseWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
