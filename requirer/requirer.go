package requirer

import "github.com/your-org/envlock/schema"

// Result holds the outcome of a required-variable check.
type Result struct {
	Key     string
	Present bool
	Value   string
}

// Report is the full output of a Require call.
type Report struct {
	Results []Result
	Missing []string
	OK      bool
}

// Require checks that every variable marked required=true in the schema
// has a non-empty value in the supplied env map. Variables that are not
// required are ignored entirely.
func Require(s schema.Schema, env map[string]string) Report {
	var results []Result
	var missing []string

	for _, v := range s.Vars {
		if !v.Required {
			continue
		}

		val, present := env[v.Key]
		if !present || val == "" {
			// Fall back to schema default before declaring missing.
			if v.Default != "" {
				results = append(results, Result{Key: v.Key, Present: true, Value: v.Default})
				continue
			}
			missing = append(missing, v.Key)
			results = append(results, Result{Key: v.Key, Present: false})
			continue
		}

		results = append(results, Result{Key: v.Key, Present: true, Value: val})
	}

	return Report{
		Results: results,
		Missing: missing,
		OK:      len(missing) == 0,
	}
}
