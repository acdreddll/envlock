// Package trimmer removes unused or redundant entries from a schema.
// An entry is considered unused if it has no description, no default,
// is not required, not sensitive, and has no tags or group.
package trimmer

import (
	"github.com/your-org/envlock/schema"
)

// Result holds the outcome of a trim operation.
type Result struct {
	Removed []string
	Kept    []schema.EnvVar
}

// TrimMode controls which entries are eligible for removal.
type TrimMode int

const (
	// TrimBare removes entries that are completely unannotated.
	TrimBare TrimMode = iota
	// TrimOptionalNoDefault removes optional entries that have no default value.
	TrimOptionalNoDefault
)

// Trim scans the schema and removes entries matching the given mode.
// It returns a Result describing what was removed and what was kept.
func Trim(s schema.Schema, mode TrimMode) Result {
	var result Result

	for _, ev := range s.Vars {
		switch mode {
		case TrimBare:
			if isBare(ev) {
				result.Removed = append(result.Removed, ev.Key)
			} else {
				result.Kept = append(result.Kept, ev)
			}
		case TrimOptionalNoDefault:
			if !ev.Required && ev.Default == "" {
				result.Removed = append(result.Removed, ev.Key)
			} else {
				result.Kept = append(result.Kept, ev)
			}
		}
	}

	return result
}

// isBare returns true when an EnvVar carries no meaningful metadata.
func isBare(ev schema.EnvVar) bool {
	return ev.Description == "" &&
		ev.Default == "" &&
		!ev.Required &&
		!ev.Sensitive &&
		len(ev.Tags) == 0 &&
		ev.Group == ""
}
