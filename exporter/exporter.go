package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yourorg/envlock/validator"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Export writes the resolved environment variables (with defaults applied)
// to the given writer in the specified format.
func Export(w io.Writer, report *validator.Report, format Format) error {
	keys := make([]string, 0, len(report.Resolved))
	for k := range report.Resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case FormatDotEnv:
		return writeDotEnv(w, keys, report.Resolved)
	case FormatExport:
		return writeExport(w, keys, report.Resolved)
	case FormatJSON:
		return writeJSON(w, keys, report.Resolved)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func writeDotEnv(w io.Writer, keys []string, resolved map[string]string) error {
	for _, k := range keys {
		v := resolved[k]
		v = strings.ReplaceAll(v, `"`, `\"`)
		if _, err := fmt.Fprintf(w, "%s=\"%s\"\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, keys []string, resolved map[string]string) error {
	for _, k := range keys {
		v := resolved[k]
		v = strings.ReplaceAll(v, `'`, `'\''`)
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, keys []string, resolved map[string]string) error {
	if _, err := fmt.Fprint(w, "{\n"); err != nil {
		return err
	}
	for i, k := range keys {
		v := resolved[k]
		v = strings.ReplaceAll(v, `"`, `\"`)
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  \"%s\": \"%s\"%s\n", k, v, comma); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "}\n")
	return err
}
