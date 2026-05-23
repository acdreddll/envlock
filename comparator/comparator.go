package comparator

import (
	"fmt"
	"sort"

	"github.com/yourorg/envlock/schema"
)

// FieldDiff describes a difference in a single field of an EnvVar.
type FieldDiff struct {
	Field string
	Left  string
	Right string
}

// EntryDiff holds all field-level differences for a single key.
type EntryDiff struct {
	Key    string
	Fields []FieldDiff
}

// Result is the full output of a schema comparison.
type Result struct {
	OnlyInLeft  []string
	OnlyInRight []string
	Differing   []EntryDiff
}

// HasChanges returns true when any difference was found.
func (r Result) HasChanges() bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Differing) > 0
}

// Compare performs a deep field-level comparison between two schemas.
func Compare(left, right schema.Schema) Result {
	lIdx := indexByKey(left)
	rIdx := indexByKey(right)

	result := Result{}

	for k := range lIdx {
		if _, ok := rIdx[k]; !ok {
			result.OnlyInLeft = append(result.OnlyInLeft, k)
		}
	}
	for k := range rIdx {
		if _, ok := lIdx[k]; !ok {
			result.OnlyInRight = append(result.OnlyInRight, k)
		}
	}

	for k, lv := range lIdx {
		rv, ok := rIdx[k]
		if !ok {
			continue
		}
		if diffs := fieldDiffs(lv, rv); len(diffs) > 0 {
			result.Differing = append(result.Differing, EntryDiff{Key: k, Fields: diffs})
		}
	}

	sort.Strings(result.OnlyInLeft)
	sort.Strings(result.OnlyInRight)
	sort.Slice(result.Differing, func(i, j int) bool {
		return result.Differing[i].Key < result.Differing[j].Key
	})

	return result
}

func indexByKey(s schema.Schema) map[string]schema.EnvVar {
	m := make(map[string]schema.EnvVar, len(s.Vars))
	for _, v := range s.Vars {
		m[v.Key] = v
	}
	return m
}

func fieldDiffs(a, b schema.EnvVar) []FieldDiff {
	var diffs []FieldDiff
	check := func(field, l, r string) {
		if l != r {
			diffs = append(diffs, FieldDiff{Field: field, Left: l, Right: r})
		}
	}
	check("description", a.Description, b.Description)
	check("default", a.Default, b.Default)
	check("pattern", a.Pattern, b.Pattern)
	check("required", fmt.Sprintf("%v", a.Required), fmt.Sprintf("%v", b.Required))
	check("sensitive", fmt.Sprintf("%v", a.Sensitive), fmt.Sprintf("%v", b.Sensitive))
	return diffs
}
