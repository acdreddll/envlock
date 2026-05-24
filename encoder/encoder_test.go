package encoder_test

import (
	"encoding/base64"
	"testing"

	"github.com/envlock/encoder"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, sensitive bool) schema.EnvVar {
	return schema.EnvVar{Key: key, Sensitive: sensitive}
}

func TestEncode_Base64AllVars(t *testing.T) {
	s := makeSchema(ev("DB_PASS", true), ev("APP_NAME", false))
	env := map[string]string{"DB_PASS": "secret", "APP_NAME": "myapp"}

	results, err := encoder.Encode(s, env, encoder.Options{Format: encoder.FormatBase64})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	expected := base64.StdEncoding.EncodeToString([]byte("secret"))
	if results[0].Encoded != expected {
		t.Errorf("expected %q, got %q", expected, results[0].Encoded)
	}
}

func TestEncode_HexFormat(t *testing.T) {
	s := makeSchema(ev("TOKEN", true))
	env := map[string]string{"TOKEN": "ab"}

	results, err := encoder.Encode(s, env, encoder.Options{Format: encoder.FormatHex})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Encoded != "6162" {
		t.Errorf("expected '6162', got %q", results[0].Encoded)
	}
}

func TestEncode_SensitiveOnly(t *testing.T) {
	s := makeSchema(ev("DB_PASS", true), ev("APP_NAME", false))
	env := map[string]string{"DB_PASS": "secret", "APP_NAME": "myapp"}

	results, err := encoder.Encode(s, env, encoder.Options{Format: encoder.FormatBase64, SensitiveOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "APP_NAME" && !r.Skipped {
			t.Errorf("APP_NAME should be skipped in SensitiveOnly mode")
		}
		if r.Key == "DB_PASS" && r.Skipped {
			t.Errorf("DB_PASS should be encoded")
		}
	}
}

func TestEncode_MissingEnvSkipped(t *testing.T) {
	s := makeSchema(ev("MISSING_KEY", false))
	env := map[string]string{}

	results, err := encoder.Encode(s, env, encoder.Options{Format: encoder.FormatBase64})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Errorf("expected MISSING_KEY to be skipped")
	}
}

func TestEncode_UnsupportedFormat(t *testing.T) {
	s := makeSchema(ev("KEY", false))
	_, err := encoder.Encode(s, map[string]string{"KEY": "val"}, encoder.Options{Format: "rot13"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestEncode_NoneFormat(t *testing.T) {
	s := makeSchema(ev("KEY", false))
	env := map[string]string{"KEY": "plaintext"}

	results, err := encoder.Encode(s, env, encoder.Options{Format: encoder.FormatNone})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Encoded != "plaintext" {
		t.Errorf("expected passthrough, got %q", results[0].Encoded)
	}
}
