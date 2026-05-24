package linter

// Summary holds aggregate statistics for a lint run.
type Summary struct {
	Total    int
	ByKind   map[string]int
	Affected []string
}

// Summarize aggregates a slice of Issues into a Summary.
func Summarize(issues []Issue) Summary {
	s := Summary{
		Total:  len(issues),
		ByKind: make(map[string]int),
	}

	seen := map[string]struct{}{}
	for _, iss := range issues {
		s.ByKind[iss.Kind]++
		if _, ok := seen[iss.Key]; !ok {
			seen[iss.Key] = struct{}{}
			s.Affected = append(s.Affected, iss.Key)
		}
	}
	return s
}

// HasIssues returns true when the summary contains at least one issue.
func (s Summary) HasIssues() bool {
	return s.Total > 0
}

// KindsPresent returns a deduplicated list of issue kinds found.
func (s Summary) KindsPresent() []string {
	kinds := make([]string, 0, len(s.ByKind))
	for k := range s.ByKind {
		kinds = append(kinds, k)
	}
	return kinds
}
