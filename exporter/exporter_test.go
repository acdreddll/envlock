package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlock/exporter"
	"github.com/yourorg/envlock/validator"
)

func makeReport(resolved map[string]string) *validator.Report {
	return &validator.Report{
		Resolved: resolved,
	}
}

func TestExport_DotEnv(t *testing.T) {
	report := makeReport(map[string]string{
		"APP_ENV": "production",
		"DB_URL":  "postgres://localhost/db",
	})

	var buf strings.Builder
	if err := exporter.Export(&buf, report, exporter.FormatDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `APP_ENV="production"`) {
		t.Errorf("expected APP_ENV in output, got:\n%s", out)
	}
	if !strings.Contains(out, `DB_URL="postgres://localhost/db"`) {
		t.Errorf("expected DB_URL in output, got:\n%s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	report := makeReport(map[string]string{
		"PORT": "8080",
	})

	var buf strings.Builder
	if err := exporter.Export(&buf, report, exporter.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "export PORT='8080'") {
		t.Errorf("expected export PORT in output, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	report := makeReport(map[string]string{
		"LOG_LEVEL": "info",
	})

	var buf strings.Builder
	if err := exporter.Export(&buf, report, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"LOG_LEVEL": "info"`) {
		t.Errorf("expected LOG_LEVEL in JSON output, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected valid JSON braces, got:\n%s", out)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	report := makeReport(map[string]string{"X": "1"})
	var buf strings.Builder
	err := exporter.Export(&buf, report, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	report := makeReport(map[string]string{
		"Z_VAR": "last",
		"A_VAR": "first",
		"M_VAR": "middle",
	})

	var buf strings.Builder
	if err := exporter.Export(&buf, report, exporter.FormatDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_VAR") {
		t.Errorf("expected A_VAR first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_VAR") {
		t.Errorf("expected Z_VAR last, got: %s", lines[2])
	}
}

// TestExport_EmptyReport verifies that exporting an empty resolved map
// produces no output and does not return an error.
func TestExport_EmptyReport(t *testing.T) {
	formats := []exporter.Format{
		exporter.FormatDotEnv,
		exporter.FormatExport,
		exporter.FormatJSON,
	}
	for _, fmt := range formats {
		var buf strings.Builder
		if err := exporter.Export(&buf, makeReport(nil), fmt); err != nil {
			t.Errorf("format %q: unexpected error for empty report: %v", fmt, err)
		}
	}
}
