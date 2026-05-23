// Package stripper removes specified fields from schema variables,
// producing a cleaned schema suitable for sharing or publishing.
package stripper

import (
	"fmt"
	"strings"

	"github.com/your-org/envlock/schema"
)

// ValidFields are the fields that can be stripped from a variable.
var ValidFields = []string{"default", "description", "pattern", "group", "tags", "sensitive", "deprecated", "remove_by"}

// Options controls which fields are removed.
type Options struct {
	Fields []string // field names to strip
}

// Strip returns a new schema with the specified fields removed from every variable.
// Unknown field names produce an error.
func Strip(s schema.Schema, opts Options) (schema.Schema, error) {
	for _, f := range opts.Fields {
		if !isValid(f) {
			return schema.Schema{}, fmt.Errorf("unknown field %q: valid fields are %s",
				f, strings.Join(ValidFields, ", "))
		}
	}

	out := make(schema.Schema, 0, len(s))
	for _, v := range s {
		copy := v
		for _, f := range opts.Fields {
			applyStrip(&copy, f)
		}
		out = append(out, copy)
	}
	return out, nil
}

func applyStrip(v *schema.EnvVar, field string) {
	switch field {
	case "default":
		v.Default = ""
	case "description":
		v.Description = ""
	case "pattern":
		v.Pattern = ""
	case "group":
		v.Group = ""
	case "tags":
		v.Tags = nil
	case "sensitive":
		v.Sensitive = false
	case "deprecated":
		v.Deprecated = false
	case "remove_by":
		v.RemoveBy = ""
	}
}

func isValid(field string) bool {
	for _, f := range ValidFields {
		if f == field {
			return true
		}
	}
	return false
}
