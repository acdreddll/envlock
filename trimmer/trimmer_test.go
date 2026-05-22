package trimmer_test

import (
	"testing"

	"github.com/your-org/envlock/schema"
	"github.com/your-org/envlock/trimmer"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, desc, def string, required, sensitive bool, tags []string, group string) schema.EnvVar {
	return schema.EnvVar{
		Key:         key,
		Description: desc,
		Default:     def,
		Required:    required,
		Sensitive:   sensitive,
		Tags:        tags,
		Group:       group,
	}
}

func TestTrim_BareRemovesUnannotated(t *testing.T) {
	s := makeSchema(
		ev("BARE_VAR", "", "", false, false, nil, ""),
		ev("WITH_DESC", "has description", "", false, false, nil, ""),
	)
	r := trimmer.Trim(s, trimmer.TrimBare)
	if len(r.Removed) != 1 || r.Removed[0] != "BARE_VAR" {
		t.Errorf("expected BARE_VAR removed, got %v", r.Removed)
	}
	if len(r.Kept) != 1 || r.Kept[0].Key != "WITH_DESC" {
		t.Errorf("expected WITH_DESC kept, got %v", r.Kept)
	}
}

func TestTrim_BareKeepsSensitive(t *testing.T) {
	s := makeSchema(
		ev("SECRET", "", "", false, true, nil, ""),
	)
	r := trimmer.Trim(s, trimmer.TrimBare)
	if len(r.Removed) != 0 {
		t.Errorf("expected nothing removed, got %v", r.Removed)
	}
}

func TestTrim_BareKeepsTagged(t *testing.T) {
	s := makeSchema(
		ev("TAGGED", "", "", false, false, []string{"infra"}, ""),
	)
	r := trimmer.Trim(s, trimmer.TrimBare)
	if len(r.Removed) != 0 {
		t.Errorf("expected nothing removed, got %v", r.Removed)
	}
}

func TestTrim_OptionalNoDefault(t *testing.T) {
	s := makeSchema(
		ev("OPT_NO_DEF", "optional", "", false, false, nil, ""),
		ev("OPT_WITH_DEF", "optional", "fallback", false, false, nil, ""),
		ev("REQUIRED", "required", "", true, false, nil, ""),
	)
	r := trimmer.Trim(s, trimmer.TrimOptionalNoDefault)
	if len(r.Removed) != 1 || r.Removed[0] != "OPT_NO_DEF" {
		t.Errorf("expected OPT_NO_DEF removed, got %v", r.Removed)
	}
	if len(r.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(r.Kept))
	}
}

func TestTrim_EmptySchema(t *testing.T) {
	r := trimmer.Trim(makeSchema(), trimmer.TrimBare)
	if len(r.Removed) != 0 || len(r.Kept) != 0 {
		t.Error("expected empty result for empty schema")
	}
}
