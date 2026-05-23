package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/envlock/schema"
)

// Result holds the hash output for a schema or individual variable.
type Result struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

// Summary is the top-level output of a Hash operation.
type Summary struct {
	SchemaHash string   `json:"schema_hash"`
	Vars       []Result `json:"vars"`
}

// Hash computes a deterministic SHA-256 fingerprint for each variable
// in the schema and a combined fingerprint for the entire schema.
// The combined hash is derived from the sorted per-variable hashes so
// that key ordering in the source YAML does not affect the result.
func Hash(s schema.Schema) Summary {
	results := make([]Result, 0, len(s.Vars))

	for _, v := range s.Vars {
		h := varHash(v)
		results = append(results, Result{Key: v.Key, Hash: h})
	}

	// Sort by key for a stable per-var list.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	// Build combined hash from sorted individual hashes.
	combined := sha256.New()
	for _, r := range results {
		fmt.Fprintf(combined, "%s:%s\n", r.Key, r.Hash)
	}
	schemaHash := hex.EncodeToString(combined.Sum(nil))

	return Summary{
		SchemaHash: schemaHash,
		Vars:       results,
	}
}

// varHash returns a hex-encoded SHA-256 digest for a single variable,
// incorporating all meaningful fields so any change is detectable.
func varHash(v schema.EnvVar) string {
	h := sha256.New()
	fmt.Fprintf(h, "key=%s\n", v.Key)
	fmt.Fprintf(h, "description=%s\n", v.Description)
	fmt.Fprintf(h, "default=%s\n", v.Default)
	fmt.Fprintf(h, "required=%t\n", v.Required)
	fmt.Fprintf(h, "sensitive=%t\n", v.Sensitive)
	fmt.Fprintf(h, "pattern=%s\n", v.Pattern)
	fmt.Fprintf(h, "group=%s\n", v.Group)

	sortedTags := make([]string, len(v.Tags))
	copy(sortedTags, v.Tags)
	sort.Strings(sortedTags)
	for _, t := range sortedTags {
		fmt.Fprintf(h, "tag=%s\n", t)
	}

	return hex.EncodeToString(h.Sum(nil))
}
