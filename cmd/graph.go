package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envlock/grapher"
	"github.com/user/envlock/schema"
)

func init() {
	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Show dependency graph of environment variables",
		RunE:  runGraph,
	}
	graphCmd.Flags().String("schema", "envlock.yaml", "path to schema file")
	graphCmd.Flags().String("format", "text", "output format: text or json")
	graphCmd.Flags().Bool("cycles-only", false, "only report cyclic dependencies")
	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, _ []string) error {
	schemaPath, _ := cmd.Flags().GetString("schema")
	format, _ := cmd.Flags().GetString("format")
	cyclesOnly, _ := cmd.Flags().GetBool("cycles-only")

	s, err := schema.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	g := grapher.Build(s)
	cycles := g.Cycles()

	if format == "json" {
		return printGraphJSON(g, cycles, cyclesOnly)
	}
	return printGraphText(g, cycles, cyclesOnly)
}

func printGraphText(g grapher.Graph, cycles [][]string, cyclesOnly bool) error {
	if len(cycles) > 0 {
		fmt.Fprintln(os.Stderr, "⚠  Cyclic dependencies detected:")
		for _, c := range cycles {
			fmt.Fprintf(os.Stderr, "   %s\n", strings.Join(c, " → "))
		}
	}
	if cyclesOnly {
		if len(cycles) == 0 {
			fmt.Println("No cycles found.")
		}
		return nil
	}

	keys := make([]string, 0, len(g.Edges))
	for k := range g.Edges {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Dependency graph:")
	for _, k := range keys {
		deps := g.Edges[k]
		if len(deps) == 0 {
			fmt.Printf("  %s (no deps)\n", k)
		} else {
			fmt.Printf("  %s → %s\n", k, strings.Join(deps, ", "))
		}
	}
	roots := g.Roots()
	if len(roots) > 0 {
		fmt.Printf("\nRoots (not depended on): %s\n", strings.Join(roots, ", "))
	}
	return nil
}

func printGraphJSON(g grapher.Graph, cycles [][]string, cyclesOnly bool) error {
	type output struct {
		Edges      map[string][]string `json:"edges,omitempty"`
		Roots      []string            `json:"roots,omitempty"`
		Cycles     [][]string          `json:"cycles"`
		CyclesOnly bool                `json:"cycles_only"`
	}
	out := output{Cycles: cycles, CyclesOnly: cyclesOnly}
	if !cyclesOnly {
		out.Edges = g.Edges
		out.Roots = g.Roots()
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
