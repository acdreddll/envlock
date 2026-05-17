package summarizer_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/summarizer"
)

func TestSummarize_FullSchema(t *testing.T) {
	entries := []schema.EnvVar{
		{Key: "DB_HOST", Required: true, Group: "db", Tags: []string{"infra"}},
		{Key: "DB_PASS", Required: true, Sensitive: true, Pattern: `^.{8,}$`, Group: "db", Tags: []string{"infra", "secret"}},
		{Key: "APP_ENV", Required: true, Default: "production", Group: "app"},
		{Key: "LOG_LEVEL", Default: "info", Group: "app", Tags: []string{"ops"}},
		{Key: "FEATURE_FLAG"},
	}

	s := summarizer.Summarize(entries)

	if s.Total != 5 {
		t.Errorf("total: want 5, got %d", s.Total)
	}
	if s.Required != 3 {
		t.Errorf("required: want 3, got %d", s.Required)
	}
	if s.Optional != 2 {
		t.Errorf("optional: want 2, got %d", s.Optional)
	}
	if s.Sensitive != 1 {
		t.Errorf("sensitive: want 1, got %d", s.Sensitive)
	}
	if s.WithDefault != 2 {
		t.Errorf("with_default: want 2, got %d", s.WithDefault)
	}
	if s.WithPattern != 1 {
		t.Errorf("with_pattern: want 1, got %d", s.WithPattern)
	}
	if s.Groups["db"] != 2 {
		t.Errorf("groups[db]: want 2, got %d", s.Groups["db"])
	}
	if s.Groups["app"] != 2 {
		t.Errorf("groups[app]: want 2, got %d", s.Groups["app"])
	}
	if s.Groups["(ungrouped)"] != 1 {
		t.Errorf("groups[ungrouped]: want 1, got %d", s.Groups["(ungrouped)"])
	}
	if s.TagCounts["infra"] != 2 {
		t.Errorf("tags[infra]: want 2, got %d", s.TagCounts["infra"])
	}
	if s.TagCounts["secret"] != 1 {
		t.Errorf("tags[secret]: want 1, got %d", s.TagCounts["secret"])
	}
	if s.TagCounts["ops"] != 1 {
		t.Errorf("tags[ops]: want 1, got %d", s.TagCounts["ops"])
	}

	groups := summarizer.SortedGroups(s)
	if groups[0] != "(ungrouped)" {
		t.Errorf("first sorted group: want (ungrouped), got %s", groups[0])
	}
}
