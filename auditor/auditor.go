package auditor

import (
	"fmt"
	"time"

	"github.com/user/envlock/schema"
)

// Severity represents the level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Finding represents a single audit result entry.
type Finding struct {
	Key      string
	Severity Severity
	Message  string
}

// Report holds all findings from an audit run.
type Report struct {
	Timestamp time.Time
	Service   string
	Findings  []Finding
}

// HasIssues returns true if any findings exist.
func (r *Report) HasIssues() bool {
	return len(r.Findings) > 0
}

// Summary returns a short human-readable summary.
func (r *Report) Summary() string {
	if !r.HasIssues() {
		return fmt.Sprintf("[%s] audit passed: no issues found", r.Service)
	}
	return fmt.Sprintf("[%s] audit found %d issue(s)", r.Service, len(r.Findings))
}

// Audit inspects a schema and the provided env map for security and hygiene issues.
func Audit(s schema.Schema, env map[string]string) Report {
	report := Report{
		Timestamp: time.Now().UTC(),
		Service:   s.Service,
	}

	for _, v := range s.Vars {
		val, present := env[v.Key]

		// Warn about sensitive keys that appear to have weak/default values.
		if v.Sensitive && present {
			if val == "" || val == "changeme" || val == "secret" || val == "password" {
				report.Findings = append(report.Findings, Finding{
					Key:      v.Key,
					Severity: SeverityCritical,
					Message:  "sensitive variable has a weak or placeholder value",
				})
			}
		}

		// Warn if a sensitive variable has no pattern constraint.
		if v.Sensitive && v.Pattern == "" {
			report.Findings = append(report.Findings, Finding{
				Key:      v.Key,
				Severity: SeverityWarning,
				Message:  "sensitive variable has no pattern constraint defined",
			})
		}

		// Info: required variable is missing from env (complements validator).
		if v.Required && !present {
			report.Findings = append(report.Findings, Finding{
				Key:      v.Key,
				Severity: SeverityCritical,
				Message:  "required variable is absent from environment",
			})
		}
	}

	return report
}
