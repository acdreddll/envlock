package auditor_test

import (
	"testing"

	"github.com/user/envlock/auditor"
	"github.com/user/envlock/schema"
)

func makeSchema(vars []schema.VarDef) schema.Schema {
	return schema.Schema{Service: "test-svc", Vars: vars}
}

func TestAudit_NoIssues(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Key: "API_KEY", Required: true, Sensitive: true, Pattern: `^[A-Za-z0-9]{32}$`},
	})
	env := map[string]string{"API_KEY": "abcdefghijklmnopqrstuvwxyz123456"}

	r := auditor.Audit(s, env)
	if r.HasIssues() {
		t.Fatalf("expected no issues, got: %+v", r.Findings)
	}
}

func TestAudit_WeakSensitiveValue(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Key: "DB_PASSWORD", Required: true, Sensitive: true, Pattern: `^.{8,}$`},
	})
	env := map[string]string{"DB_PASSWORD": "changeme"}

	r := auditor.Audit(s, env)
	if !r.HasIssues() {
		t.Fatal("expected issues for weak sensitive value")
	}
	found := false
	for _, f := range r.Findings {
		if f.Key == "DB_PASSWORD" && f.Severity == auditor.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected critical finding for DB_PASSWORD")
	}
}

func TestAudit_SensitiveNoPattern(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Key: "SECRET_TOKEN", Required: false, Sensitive: true},
	})
	env := map[string]string{"SECRET_TOKEN": "sometoken"}

	r := auditor.Audit(s, env)
	if !r.HasIssues() {
		t.Fatal("expected warning for sensitive var without pattern")
	}
	for _, f := range r.Findings {
		if f.Key == "SECRET_TOKEN" && f.Severity == auditor.SeverityWarning {
			return
		}
	}
	t.Error("expected warning finding for SECRET_TOKEN")
}

func TestAudit_RequiredMissing(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Key: "REQUIRED_VAR", Required: true},
	})
	env := map[string]string{}

	r := auditor.Audit(s, env)
	if !r.HasIssues() {
		t.Fatal("expected critical finding for missing required var")
	}
	for _, f := range r.Findings {
		if f.Key == "REQUIRED_VAR" && f.Severity == auditor.SeverityCritical {
			return
		}
	}
	t.Error("expected critical finding for REQUIRED_VAR")
}

func TestReport_Summary(t *testing.T) {
	s := makeSchema([]schema.VarDef{})
	r := auditor.Audit(s, map[string]string{})
	if r.HasIssues() {
		t.Fatal("expected clean report")
	}
	summary := r.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
