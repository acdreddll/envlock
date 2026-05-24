package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envlock/schema"
	"github.com/envlock/snapshotter"
	"github.com/spf13/cobra"
)

var snapshotLabel string
var snapshotCompare string
var snapshotFormat string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Take or compare schema snapshots",
	}

	takeCmd := &cobra.Command{
		Use:   "take",
		Short: "Capture a snapshot of the current schema",
		RunE:  runSnapshotTake,
	}
	takeCmd.Flags().StringVar(&snapshotLabel, "label", "", "Label for the snapshot")
	takeCmd.Flags().StringVar(&snapshotFormat, "format", "json", "Output format: json")

	diffCmd := &cobra.Command{
		Use:   "diff <before.json> <after.json>",
		Short: "Show drift between two snapshots",
		Args:  cobra.ExactArgs(2),
		RunE:  runSnapshotDiff,
	}
	diffCmd.Flags().StringVar(&snapshotFormat, "format", "text", "Output format: text|json")

	snapshotCmd.AddCommand(takeCmd, diffCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotTake(cmd *cobra.Command, args []string) error {
	s, err := schema.Load(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	snap := snapshotter.Take(s, snapshotLabel)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

func runSnapshotDiff(cmd *cobra.Command, args []string) error {
	before, err := loadSnapshot(args[0])
	if err != nil {
		return fmt.Errorf("load before snapshot: %w", err)
	}
	after, err := loadSnapshot(args[1])
	if err != nil {
		return fmt.Errorf("load after snapshot: %w", err)
	}
	entries := snapshotter.Diff(before, after)
	if snapshotFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}
	if len(entries) == 0 {
		fmt.Println("No drift detected.")
		return nil
	}
	for _, e := range entries {
		if e.From != "" || e.To != "" {
			fmt.Printf("  %-30s %-14s %s -> %s\n", e.Key, e.Change, e.From, e.To)
		} else {
			fmt.Printf("  %-30s %s\n", e.Key, e.Change)
		}
	}
	return nil
}

func loadSnapshot(path string) (snapshotter.Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return snapshotter.Snapshot{}, err
	}
	defer f.Close()
	var snap snapshotter.Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return snapshotter.Snapshot{}, err
	}
	return snap, nil
}
