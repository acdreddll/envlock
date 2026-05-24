package inspector_test

import (
	"testing"

	"github.com/your-org/envlock/inspector"
	"github.com/your-org/envlock/schema"
)

func TestInspect_FullWorkflow(t *testing.T) {
	s := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "API_KEY", Sensitive: true, Required: true, Tags: []string{"auth"}, Group: "api"},
			{Key: "API_URL", Required: true, Group: "api", Tags: []string{"core"}},
			{Key: "CACHE_TTL", Default: "300", Group: "cache"},
			{Key: "OLD_FLAG", Deprecated: true, RemoveBy: "v3.0"},
		},
	}

	r := inspector.Inspect(s)

	if r.Total != 4 {
		t.Fatalf("expected total 4, got %d", r.Total)
	}

	apiKeys := r.ByGroup["api"]
	if len(apiKeys) != 2 {
		t.Errorf("expected 2 api group vars, got %d", len(apiKeys))
	}

	authTagged := r.ByTag["auth"]
	if len(authTagged) != 1 || authTagged[0] != "API_KEY" {
		t.Errorf("unexpected auth tag result: %v", authTagged)
	}

	info, ok := inspector.Find(s, "OLD_FLAG")
	if !ok {
		t.Fatal("expected to find OLD_FLAG")
	}
	if !info.Deprecated {
		t.Error("expected OLD_FLAG deprecated")
	}
	if info.RemoveBy != "v3.0" {
		t.Errorf("expected remove_by v3.0, got %s", info.RemoveBy)
	}

	cache, ok := inspector.Find(s, "CACHE_TTL")
	if !ok {
		t.Fatal("expected to find CACHE_TTL")
	}
	if !cache.HasDefault || cache.Default != "300" {
		t.Errorf("expected default 300, got %q", cache.Default)
	}
}
