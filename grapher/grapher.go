// Package grapher builds a dependency graph of environment variables
// based on interpolation references (e.g. ${OTHER_VAR}).
package grapher

import (
	"regexp"
	"sort"

	"github.com/user/envlock/schema"
)

// refPattern matches ${VAR_NAME} style references.
var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Graph represents directed edges: key -> keys it depends on.
type Graph struct {
	Edges map[string][]string
}

// Build constructs a dependency graph from the schema by scanning
// default values and descriptions for interpolation references.
func Build(s schema.Schema) Graph {
	g := Graph{Edges: make(map[string][]string)}
	for _, ev := range s {
		var deps []string
		seen := map[string]bool{}
		for _, src := range []string{ev.Default, ev.Description} {
			for _, m := range refPattern.FindAllStringSubmatch(src, -1) {
				ref := m[1]
				if ref != ev.Key && !seen[ref] {
					deps = append(deps, ref)
					seen[ref] = true
				}
			}
		}
		sort.Strings(deps)
		g.Edges[ev.Key] = deps
	}
	return g
}

// Roots returns keys that have no incoming edges (nothing depends on them).
func (g Graph) Roots() []string {
	dependedOn := map[string]bool{}
	for _, deps := range g.Edges {
		for _, d := range deps {
			dependedOn[d] = true
		}
	}
	var roots []string
	for k := range g.Edges {
		if !dependedOn[k] {
			roots = append(roots, k)
		}
	}
	sort.Strings(roots)
	return roots
}

// Cycles detects any circular dependencies and returns the offending key sets.
func (g Graph) Cycles() [][]string {
	visited := map[string]bool{}
	inStack := map[string]bool{}
	var result [][]string

	var dfs func(key string, path []string)
	dfs = func(key string, path []string) {
		visited[key] = true
		inStack[key] = true
		path = append(path, key)
		for _, dep := range g.Edges[key] {
			if !visited[dep] {
				dfs(dep, path)
			} else if inStack[dep] {
				// record cycle
				cycle := make([]string, len(path))
				copy(cycle, path)
				result = append(result, cycle)
			}
		}
		inStack[key] = false
	}

	keys := make([]string, 0, len(g.Edges))
	for k := range g.Edges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !visited[k] {
			dfs(k, nil)
		}
	}
	return result
}
