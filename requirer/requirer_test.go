package requirer_test

import (
	"testing"

	"github.com/your-org/envlock/requirer"
	"github.com/your-org/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, required bool, def string) schema.EnvVar {
	return schema.EnvVar{Key: key, Required: required, Default: def}
}

func TestRequire_AllPresent(t *testing.T) {
	s := makeSchema(ev("HOST", true, ""), ev("PORT", true, ""))
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}

	r := requirer.Require(s, env)

	if !r.OK {
		t.Fatalf("expected OK, got missing: %v", r.Missing)
	}
	if len(r.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(r.Results))
	}
}

func TestRequire_MissingRequired(t *testing.T) {
	s := makeSchema(ev("DB_URL", true, ""), ev("API_KEY", true, ""))
	env := map[string]string{"DB_URL": "postgres://localhost/db"}

	r := requirer.Require(s, env)

	if r.OK {
		t.Fatal("expected not OK")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "API_KEY" {
		t.Fatalf("expected missing API_KEY, got %v", r.Missing)
	}
}

func TestRequire_OptionalSkipped(t *testing.T) {
	s := makeSchema(ev("OPTIONAL", false, ""), ev("REQUIRED", true, ""))
	env := map[string]string{"REQUIRED": "yes"}

	r := requirer.Require(s, env)

	if !r.OK {
		t.Fatalf("expected OK, got missing: %v", r.Missing)
	}
	// Only the required var should appear in Results.
	if len(r.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(r.Results))
	}
}

func TestRequire_DefaultSatisfiesRequired(t *testing.T) {
	s := makeSchema(ev("LOG_LEVEL", true, "info"))
	env := map[string]string{} // not set in env

	r := requirer.Require(s, env)

	if !r.OK {
		t.Fatalf("expected OK via default, got missing: %v", r.Missing)
	}
	if len(r.Results) != 1 || r.Results[0].Value != "info" {
		t.Fatalf("expected default value 'info', got %+v", r.Results)
	}
}

func TestRequire_EmptySchema(t *testing.T) {
	s := makeSchema()
	env := map[string]string{"ANYTHING": "value"}

	r := requirer.Require(s, env)

	if !r.OK {
		t.Fatal("expected OK for empty schema")
	}
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
}
