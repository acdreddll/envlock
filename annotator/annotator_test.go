package annotator_test

import (
	"testing"

	"github.com/your-org/envlock/annotator"
	"github.com/your-org/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string) schema.EnvVar { return schema.EnvVar{Key: key} }

func TestAnnotate_SetsDescription(t *testing.T) {
	s := makeSchema(ev("DB_HOST"))
	out, res, err := annotator.Annotate(s, "DB_HOST", annotator.Options{Description: "Database host"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars[0].Description != "Database host" {
		t.Errorf("expected description to be set, got %q", out.Vars[0].Description)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "description" {
		t.Errorf("expected applied=[description], got %v", res.Applied)
	}
}

func TestAnnotate_SkipsExistingWithoutOverwrite(t *testing.T) {
	v := ev("DB_HOST")
	v.Description = "already set"
	s := makeSchema(v)
	out, res, err := annotator.Annotate(s, "DB_HOST", annotator.Options{Description: "new value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars[0].Description != "already set" {
		t.Errorf("expected description unchanged, got %q", out.Vars[0].Description)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "description" {
		t.Errorf("expected skipped=[description], got %v", res.Skipped)
	}
}

func TestAnnotate_OverwriteReplaces(t *testing.T) {
	v := ev("API_KEY")
	v.Group = "old-group"
	s := makeSchema(v)
	out, res, err := annotator.Annotate(s, "API_KEY", annotator.Options{Group: "auth", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars[0].Group != "auth" {
		t.Errorf("expected group=auth, got %q", out.Vars[0].Group)
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied field, got %v", res.Applied)
	}
}

func TestAnnotate_KeyNotFound(t *testing.T) {
	s := makeSchema(ev("DB_HOST"))
	_, _, err := annotator.Annotate(s, "MISSING", annotator.Options{Description: "x"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAnnotate_MergesTags(t *testing.T) {
	v := ev("SVC_URL")
	v.Tags = []string{"network"}
	s := makeSchema(v)
	out, res, err := annotator.Annotate(s, "SVC_URL", annotator.Options{Tags: []string{"infra", "network"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "network" already exists; only "infra" should be added
	if len(out.Vars[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", out.Vars[0].Tags)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "tags" {
		t.Errorf("expected applied=[tags], got %v", res.Applied)
	}
}

func TestAnnotate_SetsExample(t *testing.T) {
	s := makeSchema(ev("PORT"))
	out, _, err := annotator.Annotate(s, "PORT", annotator.Options{Example: "8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars[0].Example != "8080" {
		t.Errorf("expected example=8080, got %q", out.Vars[0].Example)
	}
}
