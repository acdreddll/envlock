// Package annotator adds or updates metadata fields on schema variables
// such as description, group, tags, and examples.
package annotator

import (
	"fmt"
	"strings"

	"github.com/your-org/envlock/schema"
)

// Options controls which fields get annotated.
type Options struct {
	Description string
	Group       string
	Tags        []string
	Example     string
	Overwrite   bool // if false, skip fields already set
}

// Result describes what happened to a single variable.
type Result struct {
	Key     string
	Applied []string // list of field names that were updated
	Skipped []string // list of field names that were skipped (already set)
}

// Annotate applies Options to the named key in the schema.
// Returns a Result describing what changed, or an error if the key is not found.
func Annotate(s schema.Schema, key string, opts Options) (schema.Schema, Result, error) {
	idx := -1
	for i, v := range s.Vars {
		if v.Key == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return s, Result{}, fmt.Errorf("key %q not found in schema", key)
	}

	res := Result{Key: key}
	v := s.Vars[idx]

	if opts.Description != "" {
		if v.Description == "" || opts.Overwrite {
			v.Description = opts.Description
			res.Applied = append(res.Applied, "description")
		} else {
			res.Skipped = append(res.Skipped, "description")
		}
	}

	if opts.Group != "" {
		if v.Group == "" || opts.Overwrite {
			v.Group = opts.Group
			res.Applied = append(res.Applied, "group")
		} else {
			res.Skipped = append(res.Skipped, "group")
		}
	}

	if opts.Example != "" {
		if v.Example == "" || opts.Overwrite {
			v.Example = opts.Example
			res.Applied = append(res.Applied, "example")
		} else {
			res.Skipped = append(res.Skipped, "example")
		}
	}

	if len(opts.Tags) > 0 {
		if len(v.Tags) == 0 || opts.Overwrite {
			v.Tags = mergeTags(v.Tags, opts.Tags, opts.Overwrite)
			res.Applied = append(res.Applied, "tags")
		} else {
			res.Skipped = append(res.Skipped, "tags")
		}
	}

	s.Vars[idx] = v
	return s, res, nil
}

func mergeTags(existing, incoming []string, overwrite bool) []string {
	if overwrite {
		return incoming
	}
	seen := make(map[string]struct{}, len(existing))
	for _, t := range existing {
		seen[strings.ToLower(t)] = struct{}{}
	}
	out := append([]string{}, existing...)
	for _, t := range incoming {
		if _, ok := seen[strings.ToLower(t)]; !ok {
			out = append(out, t)
		}
	}
	return out
}
