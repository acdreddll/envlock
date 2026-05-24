package linter_test

import (
	"testing"

	"github.com/envlock/linter"
)

func issues(pairs ...string) []linter.Issue {
	out := make([]linter.Issue, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, linter.Issue{Key: pairs[i], Kind: pairs[i+1], Message: pairs[i+1]})
	}
	return out
}

func TestSummarize_Empty(t *testing.T) {
	s := linter.Summarize(nil)
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
	if s.HasIssues() {
		t.Error("expected HasIssues to be false")
	}
}

func TestSummarize_CountsByKind(t *testing.T) {
	iss := issues(
		"FOO", "missing_description",
		"BAR", "missing_description",
		"BAZ", "duplicate_key",
	)
	s := linter.Summarize(iss)

	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.ByKind["missing_description"] != 2 {
		t.Errorf("expected 2 missing_description, got %d", s.ByKind["missing_description"])
	}
	if s.ByKind["duplicate_key"] != 1 {
		t.Errorf("expected 1 duplicate_key, got %d", s.ByKind["duplicate_key"])
	}
}

func TestSummarize_AffectedKeys(t *testing.T) {
	iss := issues(
		"FOO", "missing_description",
		"FOO", "duplicate_key",
		"BAR", "missing_description",
	)
	s := linter.Summarize(iss)

	if len(s.Affected) != 2 {
		t.Errorf("expected 2 unique affected keys, got %d: %v", len(s.Affected), s.Affected)
	}
}

func TestSummarize_HasIssues(t *testing.T) {
	iss := issues("X", "duplicate_key")
	s := linter.Summarize(iss)
	if !s.HasIssues() {
		t.Error("expected HasIssues to be true")
	}
}

func TestSummarize_KindsPresent(t *testing.T) {
	iss := issues("A", "missing_description", "B", "duplicate_key")
	s := linter.Summarize(iss)
	kinds := s.KindsPresent()
	if len(kinds) != 2 {
		t.Errorf("expected 2 kinds, got %d: %v", len(kinds), kinds)
	}
}
