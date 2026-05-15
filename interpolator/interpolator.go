package interpolator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envlock/schema"
)

// Result holds the outcome of interpolating a single variable.
type Result struct {
	Key      string
	Original string
	Resolved string
	Expanded bool
}

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Interpolate expands cross-references between environment variables defined
// in the schema. A variable's default value may reference another variable
// using ${OTHER_VAR} syntax. The provided env map is used as the source of
// resolved values and is updated in-place with expanded results.
func Interpolate(s schema.Schema, env map[string]string) ([]Result, error) {
	results := make([]Result, 0, len(s.Vars))

	for _, v := range s.Vars {
		original, exists := env[v.Key]
		if !exists {
			continue
		}

		expanded, err := expand(original, env)
		if err != nil {
			return nil, fmt.Errorf("interpolating %s: %w", v.Key, err)
		}

		results = append(results, Result{
			Key:      v.Key,
			Original: original,
			Resolved: expanded,
			Expanded: expanded != original,
		})

		env[v.Key] = expanded
	}

	return results, nil
}

// expand replaces all ${VAR} references in value using the provided env map.
func expand(value string, env map[string]string) (string, error) {
	var expandErr error

	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		inner := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		resolved, ok := env[inner]
		if !ok {
			expandErr = fmt.Errorf("reference to undefined variable %q", inner)
			return match
		}
		return resolved
	})

	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
