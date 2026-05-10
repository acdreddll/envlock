package validator

import (
	"fmt"
	"os"
	"regexp"

	"github.com/envlock/schema"
)

// Result holds the outcome of a single variable validation.
type Result struct {
	Name    string
	Missing bool
	Invalid bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
	Valid   bool
}

// Validate checks the current environment against the provided schema.
func Validate(s *schema.Schema) (*Report, error) {
	if s == nil {
		return nil, fmt.Errorf("schema must not be nil")
	}

	report := &Report{Valid: true}

	for _, v := range s.Vars {
		val, exists := os.LookupEnv(v.Name)

		if !exists || val == "" {
			if v.Default != "" {
				// Use default — not a failure
				continue
			}
			if v.Required {
				report.Results = append(report.Results, Result{
					Name:    v.Name,
					Missing: true,
					Message: fmt.Sprintf("%s is required but not set", v.Name),
				})
				report.Valid = false
				continue
			}
		}

		if v.Pattern != "" && val != "" {
			matched, err := regexp.MatchString(v.Pattern, val)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for %s: %w", v.Name, err)
			}
			if !matched {
				report.Results = append(report.Results, Result{
					Name:    v.Name,
					Invalid: true,
					Message: fmt.Sprintf("%s value %q does not match pattern %q", v.Name, val, v.Pattern),
				})
				report.Valid = false
			}
		}
	}

	return report, nil
}
