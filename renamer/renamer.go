package renamer

import (
	"fmt"
	"regexp"

	"github.com/user/envlock/schema"
)

// RenameResult holds the outcome of a single variable rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Applied bool
	Reason  string
}

// validKey matches environment variable naming conventions.
var validKey = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Rename renames a variable key within the schema, returning a new schema and
// a result describing what happened. The original schema is not mutated.
func Rename(s schema.Schema, oldKey, newKey string) (schema.Schema, RenameResult, error) {
	result := RenameResult{OldKey: oldKey, NewKey: newKey}

	if !validKey.MatchString(newKey) {
		return s, result, fmt.Errorf("invalid key name %q: must match [A-Z_][A-Z0-9_]*", newKey)
	}

	// Check new key does not already exist.
	for _, ev := range s.Vars {
		if ev.Key == newKey {
			result.Reason = fmt.Sprintf("key %q already exists in schema", newKey)
			return s, result, fmt.Errorf(result.Reason)
		}
	}

	updated := make([]schema.EnvVar, len(s.Vars))
	copy(updated, s.Vars)

	found := false
	for i, ev := range updated {
		if ev.Key == oldKey {
			updated[i].Key = newKey
			found = true
			break
		}
	}

	if !found {
		result.Reason = fmt.Sprintf("key %q not found in schema", oldKey)
		return s, result, fmt.Errorf(result.Reason)
	}

	result.Applied = true
	result.Reason = fmt.Sprintf("renamed %q to %q", oldKey, newKey)

	newSchema := schema.Schema{Vars: updated}
	return newSchema, result, nil
}
