package auditor_test

import (
	"testing"

	"github.com/user/envlock/auditor"
	"github.com/user/envlock/schema"
)

// TestAudit_MultipleFindings ensures multiple issues are all collected.
func TestAudit_MultipleFindings(t *testing.T) {
	s := schema.Schema{
		Service: "multi-svc",
		Vars: []schema.VarDef{
			{Key: "DB_PASS", Required: true, Sensitive: true}, // no pattern + weak value
			{Key: "API_SECRET", Required: true, Sensitive: false}, // missing from env
		},
	}
	env := map[string]string{
		"DB_PASS": "password",
	}

	r := auditor.Audit(s, env)

	if !r.HasIssues() {
		t.Fatal("expected multiple findings")
	}

	severityCounts := map[auditor.Severity]int{}
	for _, f := range r.Findings {
		severityCounts[f.Severity]++
	}

	// DB_PASS: weak value (critical) + no pattern (warning) = 2
	// API_SECRET: missing required (critical) = 1
	if severityCounts[auditor.SeverityCritical] < 2 {
		t.Errorf("expected at least 2 critical findings, got %d", severityCounts[auditor.SeverityCritical])
	}
	if severityCounts[auditor.SeverityWarning] < 1 {
		t.Errorf("expected at least 1 warning finding, got %d", severityCounts[auditor.SeverityWarning])
	}
}

// TestAudit_EmptyEnvCleanSchema ensures empty env with no required vars passes.
func TestAudit_EmptyEnvCleanSchema(t *testing.T) {
	s := schema.Schema{
		Service: "empty-svc",
		Vars: []schema.VarDef{
			{Key: "OPTIONAL_VAR", Required: false, Sensitive: false},
		},
	}

	r := auditor.Audit(s, map[string]string{})
	if r.HasIssues() {
		t.Fatalf("expected no issues, got: %+v", r.Findings)
	}
}
