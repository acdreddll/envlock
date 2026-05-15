package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envlock/schema"
)

// Resolution holds the resolved value for an environment variable.
type Resolution struct {
	Key      string
	Value    string
	Source   string // "env", "default", or "missing"
	Required bool
}

// Result contains all resolutions and any errors encountered.
type Result struct {
	Resolutions []Resolution
	Errors      []string
}

// HasErrors returns true if any required variables could not be resolved.
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// Resolve attempts to resolve each variable in the schema against the provided
// environment map, falling back to defaults where available.
func Resolve(s *schema.Schema, env map[string]string) *Result {
	result := &Result{}

	for _, ev := range s.Vars {
		res := Resolution{
			Key:      ev.Key,
			Required: ev.Required,
		}

		if val, ok := env[ev.Key]; ok && strings.TrimSpace(val) != "" {
			res.Value = val
			res.Source = "env"
		} else if ev.Default != "" {
			res.Value = ev.Default
			res.Source = "default"
		} else if ev.Required {
			res.Source = "missing"
			result.Errors = append(result.Errors, fmt.Sprintf("required variable %q is not set and has no default", ev.Key))
		} else {
			res.Source = "missing"
		}

		result.Resolutions = append(result.Resolutions, res)
	}

	return result
}

// ResolveFromOS resolves schema variables against the current OS environment.
func ResolveFromOS(s *schema.Schema) *Result {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return Resolve(s, env)
}
