package encoder_test

import (
	"encoding/base64"
	"testing"

	"github.com/envlock/encoder"
)

func TestEncode_FullSchema(t *testing.T) {
	s := makeSchema(
		ev("DB_PASSWORD", true),
		ev("API_KEY", true),
		ev("LOG_LEVEL", false),
		ev("PORT", false),
	)
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"LOG_LEVEL":   "info",
		// PORT intentionally missing
	}

	results, err := encoder.Encode(s, env, encoder.Options{
		Format:        encoder.FormatBase64,
		SensitiveOnly: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	index := map[string]encoder.Result{}
	for _, r := range results {
		index[r.Key] = r
	}

	if index["PORT"].Skipped != true {
		t.Error("PORT should be skipped (missing from env)")
	}
	if index["LOG_LEVEL"].Encoded != base64.StdEncoding.EncodeToString([]byte("info")) {
		t.Errorf("LOG_LEVEL encoded mismatch")
	}
	if index["DB_PASSWORD"].Encoded != base64.StdEncoding.EncodeToString([]byte("hunter2")) {
		t.Errorf("DB_PASSWORD encoded mismatch")
	}
}

func TestEncode_SensitiveOnlyFullSchema(t *testing.T) {
	s := makeSchema(
		ev("SECRET", true),
		ev("PUBLIC", false),
	)
	env := map[string]string{"SECRET": "topsecret", "PUBLIC": "visible"}

	results, err := encoder.Encode(s, env, encoder.Options{
		Format:        encoder.FormatHex,
		SensitiveOnly: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range results {
		if r.Key == "PUBLIC" && r.Encoded != "visible" {
			t.Errorf("PUBLIC should pass through unchanged, got %q", r.Encoded)
		}
		if r.Key == "SECRET" && r.Encoded == "topsecret" {
			t.Errorf("SECRET should be hex-encoded")
		}
	}
}

// TestEncode_MissingSensitiveKey verifies that encoding returns an error when
// a sensitive (required) key is absent from the provided environment map.
func TestEncode_MissingSensitiveKey(t *testing.T) {
	s := makeSchema(
		ev("REQUIRED_SECRET", true),
	)
	env := map[string]string{} // REQUIRED_SECRET intentionally missing

	_, err := encoder.Encode(s, env, encoder.Options{
		Format:        encoder.FormatBase64,
		SensitiveOnly: false,
	})
	if err == nil {
		t.Fatal("expected an error for missing sensitive key, got nil")
	}
}
