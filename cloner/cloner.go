package cloner

import (
	"fmt"
	"regexp"

	"github.com/envlock/schema"
)

// Result holds the outcome of a clone operation.
type Result struct {
	SourceKey string
	NewKey    string
	Cloned    bool
	Reason    string
}

var validKey = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Clone duplicates an existing schema entry under a new key.
// The new entry is an exact copy of the source, with the key replaced.
// Optional overrides for description and default value may be supplied.
func Clone(s schema.Schema, sourceKey, newKey, description, defaultVal string) (schema.Schema, Result) {
	result := Result{SourceKey: sourceKey, NewKey: newKey}

	if !validKey.MatchString(newKey) {
		result.Reason = fmt.Sprintf("invalid key name %q: must match [A-Z][A-Z0-9_]*", newKey)
		return s, result
	}

	var source *schema.EnvVar
	for i := range s.Vars {
		if s.Vars[i].Key == sourceKey {
			source = &s.Vars[i]
			break
		}
	}
	if source == nil {
		result.Reason = fmt.Sprintf("source key %q not found", sourceKey)
		return s, result
	}

	for _, v := range s.Vars {
		if v.Key == newKey {
			result.Reason = fmt.Sprintf("key %q already exists", newKey)
			return s, result
		}
	}

	cloned := *source
	cloned.Key = newKey
	if description != "" {
		cloned.Description = description
	}
	if defaultVal != "" {
		cloned.Default = defaultVal
	}

	out := schema.Schema{Vars: make([]schema.EnvVar, len(s.Vars)+1)}
	copy(out.Vars, s.Vars)
	out.Vars[len(s.Vars)] = cloned

	result.Cloned = true
	return out, result
}
