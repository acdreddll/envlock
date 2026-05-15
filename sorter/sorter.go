package sorter

import (
	"sort"

	"github.com/envlock/schema"
)

// SortOrder defines the ordering strategy for schema variables.
type SortOrder string

const (
	SortByKey         SortOrder = "key"
	SortByRequired    SortOrder = "required"
	SortByDescription SortOrder = "description"
)

// Result holds the sorted list of environment variable entries.
type Result struct {
	Vars  []schema.EnvVar
	Order SortOrder
}

// Sort returns a new Result with the schema's variables sorted by the given order.
// The original schema is not mutated.
func Sort(s schema.Schema, order SortOrder) Result {
	vars := make([]schema.EnvVar, len(s.Vars))
	copy(vars, s.Vars)

	switch order {
	case SortByRequired:
		sort.SliceStable(vars, func(i, j int) bool {
			if vars[i].Required == vars[j].Required {
				return vars[i].Key < vars[j].Key
			}
			// required vars come first
			return vars[i].Required && !vars[j].Required
		})
	case SortByDescription:
		sort.SliceStable(vars, func(i, j int) bool {
			if vars[i].Description == vars[j].Description {
				return vars[i].Key < vars[j].Key
			}
			return vars[i].Description < vars[j].Description
		})
	default: // SortByKey
		sort.SliceStable(vars, func(i, j int) bool {
			return vars[i].Key < vars[j].Key
		})
	}

	return Result{Vars: vars, Order: order}
}
