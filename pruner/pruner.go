package pruner

import "github.com/your-org/envlock/schema"

// Result holds the outcome of a prune operation.
type Result struct {
	Kept    []schema.EnvVar
	Removed []schema.EnvVar
}

// Options controls what pruner considers "unused" or removable.
type Options struct {
	// RemoveDeprecated removes entries that have a non-empty Deprecated field.
	RemoveDeprecated bool
	// RemoveOptionalNoDefault removes optional vars that have no default and no description.
	RemoveOptionalNoDefault bool
	// Keys is an explicit list of keys to remove; takes precedence over other options.
	Keys []string
}

// Prune removes entries from the schema according to the given Options and
// returns a Result describing what was kept and what was removed.
func Prune(s schema.Schema, opts Options) Result {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var kept, removed []schema.EnvVar

	for _, ev := range s.Vars {
		if shouldRemove(ev, opts, keySet) {
			removed = append(removed, ev)
		} else {
			kept = append(kept, ev)
		}
	}

	return Result{Kept: kept, Removed: removed}
}

func shouldRemove(ev schema.EnvVar, opts Options, keySet map[string]struct{}) bool {
	if _, ok := keySet[ev.Key]; ok {
		return true
	}
	if opts.RemoveDeprecated && ev.Deprecated != "" {
		return true
	}
	if opts.RemoveOptionalNoDefault && !ev.Required && ev.Default == "" && ev.Description == "" {
		return true
	}
	return false
}
