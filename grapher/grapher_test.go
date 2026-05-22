package grapher_test

import (
	"testing"

	"github.com/user/envlock/grapher"
	"github.com/user/envlock/schema"
)

func makeSchema(entries ...schema.EnvVar) schema.Schema {
	return schema.Schema(entries)
}

func ev(key, def, desc string) schema.EnvVar {
	return schema.EnvVar{Key: key, Default: def, Description: desc}
}

func TestBuild_NoDependencies(t *testing.T) {
	s := makeSchema(ev("HOST", "localhost", "server host"), ev("PORT", "8080", "server port"))
	g := grapher.Build(s)
	if len(g.Edges["HOST"]) != 0 {
		t.Errorf("expected no deps for HOST, got %v", g.Edges["HOST"])
	}
}

func TestBuild_SingleDependency(t *testing.T) {
	s := makeSchema(
		ev("BASE_URL", "http://localhost", ""),
		ev("API_URL", "${BASE_URL}/api", "endpoint"),
	)
	g := grapher.Build(s)
	deps := g.Edges["API_URL"]
	if len(deps) != 1 || deps[0] != "BASE_URL" {
		t.Errorf("expected [BASE_URL], got %v", deps)
	}
	if len(g.Edges["BASE_URL"]) != 0 {
		t.Errorf("BASE_URL should have no deps")
	}
}

func TestBuild_MultipleDependencies(t *testing.T) {
	s := makeSchema(
		ev("HOST", "localhost", ""),
		ev("PORT", "8080", ""),
		ev("ADDR", "${HOST}:${PORT}", "full address"),
	)
	g := grapher.Build(s)
	deps := g.Edges["ADDR"]
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(deps), deps)
	}
	if deps[0] != "HOST" || deps[1] != "PORT" {
		t.Errorf("unexpected deps: %v", deps)
	}
}

func TestRoots_IdentifiesLeaves(t *testing.T) {
	s := makeSchema(
		ev("BASE", "http://example.com", ""),
		ev("URL", "${BASE}/path", ""),
	)
	g := grapher.Build(s)
	roots := g.Roots()
	if len(roots) != 1 || roots[0] != "URL" {
		t.Errorf("expected [URL] as root, got %v", roots)
	}
}

func TestCycles_NoCycle(t *testing.T) {
	s := makeSchema(
		ev("A", "", ""),
		ev("B", "${A}", ""),
	)
	g := grapher.Build(s)
	if cycles := g.Cycles(); len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}

func TestCycles_DetectsCycle(t *testing.T) {
	// Manually inject a cycle since schema defaults can't truly be circular
	g := grapher.Graph{
		Edges: map[string][]string{
			"A": {"B"},
			"B": {"C"},
			"C": {"A"},
		},
	}
	cycles := g.Cycles()
	if len(cycles) == 0 {
		t.Error("expected at least one cycle to be detected")
	}
}

func TestBuild_DeduplicatesRefs(t *testing.T) {
	s := makeSchema(ev("X", "${BASE}", "also uses ${BASE}"))
	g := grapher.Build(s)
	if len(g.Edges["X"]) != 1 {
		t.Errorf("expected 1 unique dep, got %v", g.Edges["X"])
	}
}
