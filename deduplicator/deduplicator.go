package deduplicator

import (
	"fmt"

	"github.com/envlock/schema"
)

// Duplicate describes a key that appears more than once in a schema.
type Duplicate struct {
	Key   string
	Count int
}

// Result holds the outcome of a deduplication pass.
type Result struct {
	Duplicates []Duplicate
	Cleaned    schema.Schema
}

// Deduplicate scans the schema for entries with identical keys and returns a
// Result containing the list of duplicates found and a cleaned schema that
// retains only the first occurrence of each key.
func Deduplicate(s schema.Schema) Result {
	seen := make(map[string]int)
	for _, v := range s.Vars {
		seen[v.Key]++
	}

	var dups []Duplicate
	for key, count := range seen {
		if count > 1 {
			dups = append(dups, Duplicate{Key: key, Count: count})
		}
	}

	// Sort duplicates by key for deterministic output.
	sortDuplicates(dups)

	cleaned := dedupe(s)
	return Result{Duplicates: dups, Cleaned: cleaned}
}

// dedupe returns a new schema keeping only the first occurrence of each key.
func dedupe(s schema.Schema) schema.Schema {
	visited := make(map[string]bool)
	var vars []schema.EnvVar
	for _, v := range s.Vars {
		if !visited[v.Key] {
			visited[v.Key] = true
			vars = append(vars, v)
		}
	}
	return schema.Schema{Vars: vars}
}

// sortDuplicates sorts a slice of Duplicate by Key alphabetically.
func sortDuplicates(dups []Duplicate) {
	for i := 1; i < len(dups); i++ {
		for j := i; j > 0 && dups[j].Key < dups[j-1].Key; j-- {
			dups[j], dups[j-1] = dups[j-1], dups[j]
		}
	}
}

// Summary returns a human-readable summary of the deduplication result.
func Summary(r Result) string {
	if len(r.Duplicates) == 0 {
		return "no duplicate keys found"
	}
	out := fmt.Sprintf("%d duplicate key(s) found:\n", len(r.Duplicates))
	for _, d := range r.Duplicates {
		out += fmt.Sprintf("  %s (%d occurrences)\n", d.Key, d.Count)
	}
	return out
}
