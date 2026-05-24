package splitter

import (
	"fmt"
	"strings"

	"github.com/envlock/schema"
)

// Result holds the two schema halves produced by Split.
type Result struct {
	Matched  schema.Schema
	Remainder schema.Schema
}

// ByGroup splits a schema into variables that belong to the given group and
// everything else.
func ByGroup(s schema.Schema, group string) Result {
	var matched, remainder schema.Schema
	for _, v := range s.Vars {
		if strings.EqualFold(v.Group, group) {
			matched.Vars = append(matched.Vars, v)
		} else {
			remainder.Vars = append(remainder.Vars, v)
		}
	}
	return Result{Matched: matched, Remainder: remainder}
}

// ByRequired splits a schema into required variables and optional ones.
func ByRequired(s schema.Schema) Result {
	var matched, remainder schema.Schema
	for _, v := range s.Vars {
		if v.Required {
			matched.Vars = append(matched.Vars, v)
		} else {
			remainder.Vars = append(remainder.Vars, v)
		}
	}
	return Result{Matched: matched, Remainder: remainder}
}

// BySensitive splits a schema into sensitive variables and non-sensitive ones.
func BySensitive(s schema.Schema) Result {
	var matched, remainder schema.Schema
	for _, v := range s.Vars {
		if v.Sensitive {
			matched.Vars = append(matched.Vars, v)
		} else {
			remainder.Vars = append(remainder.Vars, v)
		}
	}
	return Result{Matched: matched, Remainder: remainder}
}

// ByTag splits a schema into variables that carry the given tag and those that
// do not.
func ByTag(s schema.Schema, tag string) Result {
	var matched, remainder schema.Schema
	for _, v := range s.Vars {
		if hasTag(v.Tags, tag) {
			matched.Vars = append(matched.Vars, v)
		} else {
			remainder.Vars = append(remainder.Vars, v)
		}
	}
	return Result{Matched: matched, Remainder: remainder}
}

func hasTag(tags []string, target string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, target) {
			return true
		}
	}
	return false
}

// Summary returns a human-readable description of a Result.
func Summary(r Result) string {
	return fmt.Sprintf("matched=%d remainder=%d", len(r.Matched.Vars), len(r.Remainder.Vars))
}
