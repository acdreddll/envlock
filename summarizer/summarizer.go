package summarizer

import (
	"sort"

	"github.com/envlock/schema"
)

// Summary holds aggregated statistics about a schema.
type Summary struct {
	Total       int            `json:"total"`
	Required    int            `json:"required"`
	Optional    int            `json:"optional"`
	Sensitive   int            `json:"sensitive"`
	WithDefault int            `json:"with_default"`
	WithPattern int            `json:"with_pattern"`
	Groups      map[string]int `json:"groups"`
	TagCounts   map[string]int `json:"tag_counts"`
}

// Summarize computes a Summary from the provided schema entries.
func Summarize(entries []schema.EnvVar) Summary {
	s := Summary{
		Groups:    make(map[string]int),
		TagCounts: make(map[string]int),
	}

	for _, ev := range entries {
		s.Total++
		if ev.Required {
			s.Required++
		} else {
			s.Optional++
		}
		if ev.Sensitive {
			s.Sensitive++
		}
		if ev.Default != "" {
			s.WithDefault++
		}
		if ev.Pattern != "" {
			s.WithPattern++
		}
		group := ev.Group
		if group == "" {
			group = "(ungrouped)"
		}
		s.Groups[group]++
		for _, tag := range ev.Tags {
			s.TagCounts[tag]++
		}
	}
	return s
}

// SortedGroups returns group names from the summary in sorted order.
func SortedGroups(s Summary) []string {
	keys := make([]string, 0, len(s.Groups))
	for k := range s.Groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedTags returns tag names from the summary in sorted order.
func SortedTags(s Summary) []string {
	keys := make([]string, 0, len(s.TagCounts))
	for k := range s.TagCounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
