package deprecator

import (
	"fmt"
	"time"

	"github.com/envlock/schema"
)

// Finding represents a single deprecation finding for a variable.
type Finding struct {
	Key     string
	Message string
	Since   string
	RemoveBy string
}

// Result holds all findings from a deprecation check.
type Result struct {
	Findings []Finding
}

// HasIssues returns true if any deprecation findings exist.
func (r Result) HasIssues() bool {
	return len(r.Findings) > 0
}

// Deprecate inspects a schema for variables marked as deprecated and returns
// a Result describing each one. A variable is considered deprecated when its
// Deprecated field is non-empty. If RemoveBy is set and the date has passed,
// the finding message reflects that the deadline has been exceeded.
func Deprecate(s schema.Schema) Result {
	var findings []Finding
	for _, ev := range s.Vars {
		if ev.Deprecated == "" {
			continue
		}
		f := Finding{
			Key:      ev.Key,
			Since:    ev.Deprecated,
			RemoveBy: ev.RemoveBy,
		}
		if ev.RemoveBy != "" {
			deadline, err := time.Parse("2006-01-02", ev.RemoveBy)
			if err == nil && time.Now().After(deadline) {
				f.Message = fmt.Sprintf("%s is deprecated since %s and PAST its removal deadline of %s", ev.Key, ev.Deprecated, ev.RemoveBy)
			} else {
				f.Message = fmt.Sprintf("%s is deprecated since %s; scheduled for removal by %s", ev.Key, ev.Deprecated, ev.RemoveBy)
			}
		} else {
			f.Message = fmt.Sprintf("%s is deprecated since %s", ev.Key, ev.Deprecated)
		}
		findings = append(findings, f)
	}
	return Result{Findings: findings}
}
