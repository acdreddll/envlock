package profiler

import (
	"fmt"
	"sort"

	"github.com/envlock/schema"
)

// FieldProfile holds completeness metrics for a single env var entry.
type FieldProfile struct {
	Key             string
	HasDescription  bool
	HasDefault      bool
	HasPattern      bool
	HasGroup        bool
	HasTags         bool
	IsSensitive     bool
	IsRequired      bool
	CompletenessScore int // 0–100
}

// Report is the full profiling result for a schema.
type Report struct {
	Fields          []FieldProfile
	TotalFields     int
	AvgCompleteness float64
	MissingDesc     []string
	MissingDefault  []string
}

// Profile analyses every entry in the schema and returns a Report.
func Profile(s schema.Schema) Report {
	var report Report
	report.TotalFields = len(s.Vars)

	var totalScore int

	for _, v := range s.Vars {
		fp := FieldProfile{
			Key:         v.Key,
			HasDescription: v.Description != "",
			HasDefault:     v.Default != "",
			HasPattern:     v.Pattern != "",
			HasGroup:       v.Group != "",
			HasTags:        len(v.Tags) > 0,
			IsSensitive:    v.Sensitive,
			IsRequired:     v.Required,
		}

		score := 0
		if fp.HasDescription {
			score += 30
		}
		if fp.HasDefault {
			score += 20
		}
		if fp.HasPattern {
			score += 20
		}
		if fp.HasGroup {
			score += 15
		}
		if fp.HasTags {
			score += 15
		}
		fp.CompletenessScore = score
		totalScore += score

		if !fp.HasDescription {
			report.MissingDesc = append(report.MissingDesc, v.Key)
		}
		if !fp.HasDefault && !v.Required {
			report.MissingDefault = append(report.MissingDefault, v.Key)
		}

		report.Fields = append(report.Fields, fp)
	}

	sort.Slice(report.Fields, func(i, j int) bool {
		return report.Fields[i].Key < report.Fields[j].Key
	})
	sort.Strings(report.MissingDesc)
	sort.Strings(report.MissingDefault)

	if report.TotalFields > 0 {
		report.AvgCompleteness = float64(totalScore) / float64(report.TotalFields)
	}

	return report
}

// Summary returns a human-readable one-line summary of the report.
func Summary(r Report) string {
	return fmt.Sprintf("%d fields | avg completeness %.1f%% | missing desc: %d | missing default: %d",
		r.TotalFields, r.AvgCompleteness, len(r.MissingDesc), len(r.MissingDefault))
}
