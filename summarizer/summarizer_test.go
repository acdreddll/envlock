package summarizer_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/summarizer"
)

func makeSchema(evs ...schema.EnvVar) []schema.EnvVar {
	return evs
}

func ev(key string, opts ...func(*schema.EnvVar)) schema.EnvVar {
	e := schema.EnvVar{Key: key}
	for _, o := range opts {
		o(&e)
	}
	return e
}

func TestSummarize_Empty(t *testing.T) {
	s := summarizer.Summarize(makeSchema())
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
}

func TestSummarize_Counts(t *testing.T) {
	entries := makeSchema(
		ev("A", func(e *schema.EnvVar) { e.Required = true }),
		ev("B", func(e *schema.EnvVar) { e.Required = true; e.Sensitive = true }),
		ev("C", func(e *schema.EnvVar) { e.Default = "val" }),
		ev("D", func(e *schema.EnvVar) { e.Pattern = "^\\d+$" }),
	)
	s := summarizer.Summarize(entries)
	if s.Total != 4 {
		t.Errorf("expected total=4, got %d", s.Total)
	}
	if s.Required != 2 {
		t.Errorf("expected required=2, got %d", s.Required)
	}
	if s.Optional != 2 {
		t.Errorf("expected optional=2, got %d", s.Optional)
	}
	if s.Sensitive != 1 {
		t.Errorf("expected sensitive=1, got %d", s.Sensitive)
	}
	if s.WithDefault != 1 {
		t.Errorf("expected with_default=1, got %d", s.WithDefault)
	}
	if s.WithPattern != 1 {
		t.Errorf("expected with_pattern=1, got %d", s.WithPattern)
	}
}

func TestSummarize_Groups(t *testing.T) {
	entries := makeSchema(
		ev("A", func(e *schema.EnvVar) { e.Group = "auth" }),
		ev("B", func(e *schema.EnvVar) { e.Group = "auth" }),
		ev("C", func(e *schema.EnvVar) { e.Group = "db" }),
		ev("D"),
	)
	s := summarizer.Summarize(entries)
	if s.Groups["auth"] != 2 {
		t.Errorf("expected auth=2, got %d", s.Groups["auth"])
	}
	if s.Groups["db"] != 1 {
		t.Errorf("expected db=1, got %d", s.Groups["db"])
	}
	if s.Groups["(ungrouped)"] != 1 {
		t.Errorf("expected ungrouped=1, got %d", s.Groups["(ungrouped)"])
	}
}

func TestSummarize_TagCounts(t *testing.T) {
	entries := makeSchema(
		ev("A", func(e *schema.EnvVar) { e.Tags = []string{"k8s", "ci"} }),
		ev("B", func(e *schema.EnvVar) { e.Tags = []string{"k8s"} }),
	)
	s := summarizer.Summarize(entries)
	if s.TagCounts["k8s"] != 2 {
		t.Errorf("expected k8s=2, got %d", s.TagCounts["k8s"])
	}
	if s.TagCounts["ci"] != 1 {
		t.Errorf("expected ci=1, got %d", s.TagCounts["ci"])
	}
}

func TestSortedGroups(t *testing.T) {
	s := summarizer.Summary{Groups: map[string]int{"z": 1, "a": 2, "m": 3}}
	got := summarizer.SortedGroups(s)
	if got[0] != "a" || got[1] != "m" || got[2] != "z" {
		t.Errorf("unexpected order: %v", got)
	}
}
