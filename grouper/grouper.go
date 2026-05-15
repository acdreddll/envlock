package grouper

import (
	"sort"

	"github.com/envlock/schema"
)

// Group represents a named collection of environment variable entries.
type Group struct {
	Name    string
	Entries []schema.EnvVar
}

// GroupBy defines the field to group by.
type GroupBy string

const (
	GroupByRequired GroupBy = "required"
	GroupByGroup    GroupBy = "group"
	GroupBySensitive GroupBy = "sensitive"
)

// GroupResult holds the ordered list of groups produced by GroupSchema.
type GroupResult struct {
	Groups []Group
}

// GroupSchema partitions the schema's env vars into named groups based on the given field.
// Groups are returned in a deterministic (sorted) order by group name.
func GroupSchema(s schema.Schema, by GroupBy) GroupResult {
	buckets := make(map[string][]schema.EnvVar)

	for _, ev := range s.Envs {
		key := bucketKey(ev, by)
		buckets[key] = append(buckets[key], ev)
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	groups := make([]Group, 0, len(keys))
	for _, k := range keys {
		groups = append(groups, Group{
			Name:    k,
			Entries: buckets[k],
		})
	}

	return GroupResult{Groups: groups}
}

func bucketKey(ev schema.EnvVar, by GroupBy) string {
	switch by {
	case GroupByRequired:
		if ev.Required {
			return "required"
		}
		return "optional"
	case GroupBySensitive:
		if ev.Sensitive {
			return "sensitive"
		}
		return "non-sensitive"
	case GroupByGroup:
		if ev.Group != "" {
			return ev.Group
		}
		return "ungrouped"
	default:
		return "ungrouped"
	}
}
