package validator

import "fmt"

// Summary returns a human-readable summary string of the validation report.
func Summary(r Report) string {
	total := len(r.Passed) + len(r.Errors)
	if total == 0 {
		return "no variables defined in schema"
	}
	if !r.HasErrors() {
		return fmt.Sprintf("all %d variable(s) passed validation", total)
	}
	return fmt.Sprintf("%d/%d variable(s) passed, %d error(s) found", len(r.Passed), total, len(r.Errors))
}

// ErrorKeys returns the list of keys that failed validation.
func ErrorKeys(r Report) []string {
	keys := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		keys = append(keys, e.Key)
	}
	return keys
}

// PassedKeys returns the list of keys that passed validation.
func PassedKeys(r Report) []string {
	keys := make([]string, 0, len(r.Passed))
	for _, p := range r.Passed {
		keys = append(keys, p.Key)
	}
	return keys
}
