package linter

import (
	"fmt"
	"strings"

	"github.com/envlock/schema"
)

// Issue represents a linting problem found in the schema.
type Issue struct {
	Field   string
	Message string
	Severity string // "error" or "warning"
}

// Result holds the outcome of a lint run.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the lint result.
func (r *Result) Summary() string {
	if len(r.Issues) == 0 {
		return "schema is valid: no issues found"
	}
	var sb strings.Builder
	for _, issue := range r.Issues {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(issue.Severity), issue.Field, issue.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Lint checks a schema for structural and logical issues.
func Lint(s *schema.Schema) *Result {
	result := &Result{}

	keys := make(map[string]bool)
	for _, entry := range s.Vars {
		// Check for duplicate keys
		if keys[entry.Key] {
			result.Issues = append(result.Issues, Issue{
				Field:    entry.Key,
				Message:  "duplicate key defined in schema",
				Severity: "error",
			})
		}
		keys[entry.Key] = true

		// Check for empty key
		if strings.TrimSpace(entry.Key) == "" {
			result.Issues = append(result.Issues, Issue{
				Field:    "(empty)",
				Message:  "key must not be empty",
				Severity: "error",
			})
		}

		// Warn if required is true but a default is also set
		if entry.Required && entry.Default != "" {
			result.Issues = append(result.Issues, Issue{
				Field:    entry.Key,
				Message:  "required is true but a default value is set; default will never be used",
				Severity: "warning",
			})
		}

		// Warn if no description is provided
		if strings.TrimSpace(entry.Description) == "" {
			result.Issues = append(result.Issues, Issue{
				Field:    entry.Key,
				Message:  "missing description; consider documenting this variable",
				Severity: "warning",
			})
		}
	}

	return result
}
