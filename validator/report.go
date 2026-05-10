package validator

import (
	"fmt"
	"io"
	"strings"
)

// Summary returns a human-readable one-line summary of the report.
func (r *Report) Summary() string {
	if r.Valid {
		return "✓ All environment variables are valid."
	}
	return fmt.Sprintf("✗ Validation failed: %d issue(s) found.", len(r.Results))
}

// Write renders the full report to the given writer.
func (r *Report) Write(w io.Writer) {
	fmt.Fprintln(w, r.Summary())

	if len(r.Results) == 0 {
		return
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, res := range r.Results {
		kind := "INVALID"
		if res.Missing {
			kind = "MISSING"
		}
		fmt.Fprintf(w, "  [%s] %s\n", kind, res.Message)
	}
	fmt.Fprintln(w, strings.Repeat("-", 40))
}

// Errors returns all validation messages as a slice of strings.
func (r *Report) Errors() []string {
	msgs := make([]string, 0, len(r.Results))
	for _, res := range r.Results {
		msgs = append(msgs, res.Message)
	}
	return msgs
}
