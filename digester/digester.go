package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/envlock/schema"
)

// Digest holds the result of a schema digest operation.
type Digest struct {
	// SchemaHash is a deterministic hash of the entire schema structure.
	SchemaHash string `json:"schema_hash"`
	// KeyHashes maps each variable key to a hash of its definition.
	KeyHashes map[string]string `json:"key_hashes"`
	// KeyCount is the total number of variables in the schema.
	KeyCount int `json:"key_count"`
}

// Compute produces a Digest from the given schema.
// The schema hash is deterministic: same schema always produces the same hash.
func Compute(s schema.Schema) Digest {
	keyHashes := make(map[string]string, len(s.Vars))

	keys := make([]string, 0, len(s.Vars))
	for _, v := range s.Vars {
		keys = append(keys, v.Key)
	}
	sort.Strings(keys)

	index := make(map[string]schema.EnvVar, len(s.Vars))
	for _, v := range s.Vars {
		index[v.Key] = v
	}

	var sb strings.Builder
	for _, k := range keys {
		v := index[k]
		h := varDigest(v)
		keyHashes[k] = h
		fmt.Fprintf(&sb, "%s=%s\n", k, h)
	}

	overall := sha256sum(sb.String())

	return Digest{
		SchemaHash: overall,
		KeyHashes:  keyHashes,
		KeyCount:   len(s.Vars),
	}
}

// Diff returns the keys whose digests differ between two Digests.
// Returns added, removed, and changed key slices.
func Diff(a, b Digest) (added, removed, changed []string) {
	for k, bh := range b.KeyHashes {
		ah, ok := a.KeyHashes[k]
		if !ok {
			added = append(added, k)
		} else if ah != bh {
			changed = append(changed, k)
		}
	}
	for k := range a.KeyHashes {
		if _, ok := b.KeyHashes[k]; !ok {
			removed = append(removed, k)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(changed)
	return
}

func varDigest(v schema.EnvVar) string {
	parts := fmt.Sprintf(
		"key=%s|desc=%s|default=%s|required=%v|sensitive=%v|pattern=%s|group=%s|tags=%s",
		v.Key, v.Description, v.Default,
		v.Required, v.Sensitive,
		v.Pattern, v.Group,
		strings.Join(v.Tags, ","),
	)
	return sha256sum(parts)
}

func sha256sum(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
