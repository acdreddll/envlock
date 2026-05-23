package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/pinner"
	"github.com/envlock/schema"
	"github.com/spf13/cobra"
)

var pinFormat string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Capture or compare a snapshot of environment variable values",
	}

	captureCmd := &cobra.Command{
		Use:   "capture",
		Short: "Capture current env values into a snapshot (JSON to stdout)",
		RunE:  runPinCapture,
	}

	detectCmd := &cobra.Command{
		Use:   "detect [snapshot-file]",
		Short: "Detect drift between a snapshot file and current env",
		Args:  cobra.ExactArgs(1),
		RunE:  runPinDetect,
	}
	detectCmd.Flags().StringVar(&pinFormat, "format", "text", "Output format: text or json")

	pinCmd.AddCommand(captureCmd, detectCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinCapture(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	env := collectEnv(dotenvFile)
	snap := pinner.Capture(s, env)
	return json.NewEncoder(os.Stdout).Encode(snap)
}

func runPinDetect(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()

	var snap pinner.Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return fmt.Errorf("decode snapshot: %w", err)
	}
	_ = s

	env := collectEnv(dotenvFile)
	report := pinner.Detect(snap, env)

	if pinFormat == "json" {
		return json.NewEncoder(os.Stdout).Encode(report)
	}

	fmt.Println(pinner.Summary(report))
	for _, d := range report.Drifted {
		if d.Missing {
			fmt.Printf("  MISSING  %s (was %q)\n", d.Key, d.Pinned)
		} else {
			fmt.Printf("  CHANGED  %s: %q -> %q\n", d.Key, d.Pinned, d.Current)
		}
	}
	if report.HasDrift() {
		os.Exit(1)
	}
	return nil
}
