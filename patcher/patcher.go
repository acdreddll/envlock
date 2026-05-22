package patcher

import (
	"fmt"
	"regexp"

	"github.com/envlock/schema"
)

// PatchOp represents a single patch operation on a schema entry.
type PatchOp struct {
	Key         string
	Field       string
	Value       string
}

// Result holds the outcome of a patch operation.
type Result struct {
	Key     string
	Field   string
	OldVal  string
	NewVal  string
	Applied bool
	Reason  string
}

var validFields = map[string]bool{
	"description": true,
	"default":     true,
	"pattern":     true,
	"group":       true,
}

// Patch applies a list of PatchOps to a schema, returning per-op results.
func Patch(s schema.Schema, ops []PatchOp) (schema.Schema, []Result, error) {
	out := make(schema.Schema, len(s))
	copy(out, s)

	results := make([]Result, 0, len(ops))

	for _, op := range ops {
		result := Result{Key: op.Key, Field: op.Field, NewVal: op.Value}

		if !validFields[op.Field] {
			result.Applied = false
			result.Reason = fmt.Sprintf("unknown field %q", op.Field)
			results = append(results, result)
			continue
		}

		if op.Field == "pattern" && op.Value != "" {
			if _, err := regexp.Compile(op.Value); err != nil {
				result.Applied = false
				result.Reason = fmt.Sprintf("invalid pattern: %v", err)
				results = append(results, result)
				continue
			}
		}

		idx := -1
		for i, ev := range out {
			if ev.Key == op.Key {
				idx = i
				break
			}
		}

		if idx == -1 {
			result.Applied = false
			result.Reason = fmt.Sprintf("key %q not found", op.Key)
			results = append(results, result)
			continue
		}

		entry := out[idx]
		switch op.Field {
		case "description":
			result.OldVal = entry.Description
			entry.Description = op.Value
		case "default":
			result.OldVal = entry.Default
			entry.Default = op.Value
		case "pattern":
			result.OldVal = entry.Pattern
			entry.Pattern = op.Value
		case "group":
			result.OldVal = entry.Group
			entry.Group = op.Value
		}
		out[idx] = entry
		result.Applied = true
		results = append(results, result)
	}

	return out, results, nil
}
