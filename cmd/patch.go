package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/patcher"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var patchOps []string
var patchFormat string

func init() {
	patchCmd := &cobra.Command{
		Use:   "patch",
		Short: "Apply field-level patches to schema entries",
		Long:  "Patch updates specific fields (description, default, pattern, group) on schema entries without rewriting the full schema.",
		RunE:  runPatch,
	}
	patchCmd.Flags().StringArrayVarP(&patchOps, "op", "o", nil, "Patch op in KEY:FIELD:VALUE format (repeatable)")
	patchCmd.Flags().StringVarP(&patchFormat, "format", "f", "text", "Output format: text or json")
	patchCmd.Flags().StringVarP(&schemaFile, "schema", "s", "envlock.yaml", "Path to schema file")
	_ = patchCmd.MarkFlagRequired("op")
	rootCmd.AddCommand(patchCmd)
}

func runPatch(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	ops := make([]patcher.PatchOp, 0, len(patchOps))
	for _, raw := range patchOps {
		op, err := parsePatchOp(raw)
		if err != nil {
			return err
		}
		ops = append(ops, op)
	}

	_, results, err := patcher.Patch(s, ops)
	if err != nil {
		return err
	}

	if patchFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	for _, r := range results {
		if r.Applied {
			fmt.Printf("✔ %s.%s: %q → %q\n", r.Key, r.Field, r.OldVal, r.NewVal)
		} else {
			fmt.Printf("✘ %s.%s: %s\n", r.Key, r.Field, r.Reason)
		}
	}
	return nil
}

func parsePatchOp(raw string) (patcher.PatchOp, error) {
	parts := splitN(raw, ":", 3)
	if len(parts) != 3 {
		return patcher.PatchOp{}, fmt.Errorf("invalid op %q: expected KEY:FIELD:VALUE", raw)
	}
	return patcher.PatchOp{Key: parts[0], Field: parts[1], Value: parts[2]}, nil
}

func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx == -1 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
