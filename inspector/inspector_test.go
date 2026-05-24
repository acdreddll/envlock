package inspector_test

import (
	"testing"

	"github.com/your-org/envlock/inspector"
	"github.com/your-org/envlock/schema"
)

func makeSchema() schema.Schema {
	return schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_HOST", Description: "Hostname", Required: true, Group: "network", Tags: []string{"core"}},
			{Key: "DB_PASS", Description: "DB password", Sensitive: true, Group: "database", Tags: []string{"core", "secret"}},
			{Key: "LOG_LEVEL", Default: "info", Group: "logging", Deprecated: true, RemoveBy: "v2.0"},
			{Key: "RETRY_COUNT", Default: "3", Pattern: `^\d+$`},
		},
	}
}

func TestInspect_Total(t *testing.T) {
	r := inspector.Inspect(makeSchema())
	if r.Total != 4 {
		t.Fatalf("expected 4, got %d", r.Total)
	}
}

func TestInspect_ByGroup(t *testing.T) {
	r := inspector.Inspect(makeSchema())
	if len(r.ByGroup["network"]) != 1 || r.ByGroup["network"][0] != "APP_HOST" {
		t.Errorf("unexpected network group: %v", r.ByGroup["network"])
	}
	if len(r.ByGroup["database"]) != 1 {
		t.Errorf("unexpected database group: %v", r.ByGroup["database"])
	}
}

func TestInspect_ByTag(t *testing.T) {
	r := inspector.Inspect(makeSchema())
	if len(r.ByTag["core"]) != 2 {
		t.Errorf("expected 2 core-tagged vars, got %d", len(r.ByTag["core"]))
	}
	if len(r.ByTag["secret"]) != 1 {
		t.Errorf("expected 1 secret-tagged var, got %d", len(r.ByTag["secret"]))
	}
}

func TestInspect_VarDetails(t *testing.T) {
	r := inspector.Inspect(makeSchema())
	var logVar inspector.VarInfo
	for _, v := range r.Vars {
		if v.Key == "LOG_LEVEL" {
			logVar = v
		}
	}
	if !logVar.Deprecated {
		t.Error("expected LOG_LEVEL to be deprecated")
	}
	if logVar.RemoveBy != "v2.0" {
		t.Errorf("expected remove_by v2.0, got %s", logVar.RemoveBy)
	}
	if !logVar.HasDefault {
		t.Error("expected LOG_LEVEL to have a default")
	}
}

func TestFind_Found(t *testing.T) {
	info, ok := inspector.Find(makeSchema(), "DB_PASS")
	if !ok {
		t.Fatal("expected to find DB_PASS")
	}
	if !info.Sensitive {
		t.Error("expected DB_PASS to be sensitive")
	}
}

func TestFind_NotFound(t *testing.T) {
	_, ok := inspector.Find(makeSchema(), "NONEXISTENT")
	if ok {
		t.Error("expected not found")
	}
}

func TestInspect_EmptySchema(t *testing.T) {
	r := inspector.Inspect(schema.Schema{})
	if r.Total != 0 {
		t.Errorf("expected 0 total, got %d", r.Total)
	}
	if len(r.Vars) != 0 {
		t.Errorf("expected empty vars")
	}
}
