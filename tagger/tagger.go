package tagger

import (
	"fmt"
	"sort"

	"github.com/user/envlock/schema"
)

// TagIndex maps tag names to the list of env var keys that carry that tag.
type TagIndex map[string][]string

// Result holds the outcome of a tagging operation.
type Result struct {
	Index   TagIndex
	Untagged []string
}

// Tag scans a schema and builds an index of tags -> keys.
// Each EnvVar may declare zero or more tags via the Tags field.
// Keys with no tags are collected in Result.Untagged.
func Tag(s schema.Schema) (Result, error) {
	if len(s.Vars) == 0 {
		return Result{Index: make(TagIndex)}, nil
	}

	index := make(TagIndex)
	var untagged []string

	for _, ev := range s.Vars {
		if len(ev.Tags) == 0 {
			untagged = append(untagged, ev.Key)
			continue
		}
		for _, tag := range ev.Tags {
			if tag == "" {
				return Result{}, fmt.Errorf("key %q has an empty tag", ev.Key)
			}
			index[tag] = append(index[tag], ev.Key)
		}
	}

	// Sort keys within each tag bucket for deterministic output.
	for tag := range index {
		sort.Strings(index[tag])
	}
	if untagged != nil {
		sort.Strings(untagged)
	}

	return Result{Index: index, Untagged: untagged}, nil
}

// KeysForTag returns the keys associated with a given tag, or an empty slice
// if the tag is not present in the index.
func KeysForTag(r Result, tag string) []string {
	keys, ok := r.Index[tag]
	if !ok {
		return []string{}
	}
	return keys
}
