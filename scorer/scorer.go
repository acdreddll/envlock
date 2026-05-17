package scorer

import (
	"github.com/user/envlock/schema"
)

// Result holds the overall schema quality score and breakdown.
type Result struct {
	Total       int
	Score       int
	Percentage  float64
	Breakdown   map[string]int
}

// Score evaluates a schema and returns a quality score out of 100.
// Points are awarded for: descriptions, patterns on sensitive vars,
// defaults on optional vars, and group assignments.
func Score(s schema.Schema) Result {
	total := len(s.Vars)
	if total == 0 {
		return Result{Total: 0, Score: 0, Percentage: 0, Breakdown: map[string]int{}}
	}

	breakdown := map[string]int{
		"description": 0,
		"pattern":     0,
		"default":     0,
		"group":       0,
	}

	for _, v := range s.Vars {
		if v.Description != "" {
			breakdown["description"]++
		}
		if v.Sensitive && v.Pattern != "" {
			breakdown["pattern"]++
		}
		if !v.Required && v.Default != "" {
			breakdown["default"]++
		}
		if v.Group != "" {
			breakdown["group"]++
		}
	}

	// Weight: description=40, pattern=20, default=20, group=20
	weighted := float64(breakdown["description"])*40 +
		float64(breakdown["pattern"])*20 +
		float64(breakdown["default"])*20 +
		float64(breakdown["group"])*20

	max := float64(total) * 100
	pct := (weighted / max) * 100
	if pct > 100 {
		pct = 100
	}

	score := int(pct)
	return Result{
		Total:      total,
		Score:      score,
		Percentage: pct,
		Breakdown:  breakdown,
	}
}
