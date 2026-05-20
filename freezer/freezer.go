// Package freezer produces a locked snapshot of resolved environment variables
// from a schema, suitable for pinning exact values in CI or audit trails.
package freezer

import (
	"fmt"
	"sort"
	"time"

	"github.com/user/envlock/schema"
)

// Snapshot holds a point-in-time capture of resolved env values.
type Snapshot struct {
	CreatedAt time.Time            `json:"created_at" yaml:"created_at"`
	SchemaFile string              `json:"schema_file" yaml:"schema_file"`
	Entries   []SnapshotEntry      `json:"entries"    yaml:"entries"`
}

// SnapshotEntry records a single key's resolved value and metadata.
type SnapshotEntry struct {
	Key         string `json:"key"          yaml:"key"`
	Value       string `json:"value"        yaml:"value"`
	FromDefault bool   `json:"from_default" yaml:"from_default"`
	Sensitive   bool   `json:"sensitive"    yaml:"sensitive"`
}

// Freeze resolves env values against the schema and returns a Snapshot.
// env is a map of key→value (e.g. from os.Environ or a .env file).
// schemaFile is stored as metadata only.
func Freeze(s schema.Schema, env map[string]string, schemaFile string) (Snapshot, error) {
	entries := make([]SnapshotEntry, 0, len(s.Vars))

	for _, v := range s.Vars {
		val, ok := env[v.Key]
		fromDefault := false

		if !ok {
			if v.Default != "" {
				val = v.Default
				fromDefault = true
			} else if v.Required {
				return Snapshot{}, fmt.Errorf("required variable %q is not set", v.Key)
			} else {
				continue
			}
		}

		entries = append(entries, SnapshotEntry{
			Key:         v.Key,
			Value:       val,
			FromDefault: fromDefault,
			Sensitive:   v.Sensitive,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Snapshot{
		CreatedAt:  time.Now().UTC(),
		SchemaFile: schemaFile,
		Entries:    entries,
	}, nil
}

// Redacted returns a copy of the snapshot with sensitive values masked.
func (snap Snapshot) Redacted() Snapshot {
	copy := snap
	copy.Entries = make([]SnapshotEntry, len(snap.Entries))
	for i, e := range snap.Entries {
		if e.Sensitive {
			e.Value = "***"
		}
		copy.Entries[i] = e
	}
	return copy
}
