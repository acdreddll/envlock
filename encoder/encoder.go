package encoder

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/envlock/schema"
)

// Format represents the encoding format to apply.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatHex    Format = "hex"
	FormatNone   Format = "none"
)

// Result holds the encoded output for a single variable.
type Result struct {
	Key      string
	Original string
	Encoded  string
	Format   Format
	Skipped  bool
}

// Options controls which variables are encoded and how.
type Options struct {
	Format       Format
	SensitiveOnly bool
}

// Encode applies the chosen encoding to env var values.
// Only variables present in env are encoded; others are skipped.
func Encode(s schema.Schema, env map[string]string, opts Options) ([]Result, error) {
	if opts.Format == "" {
		opts.Format = FormatBase64
	}
	if opts.Format != FormatBase64 && opts.Format != FormatHex && opts.Format != FormatNone {
		return nil, fmt.Errorf("unsupported encoding format: %q", opts.Format)
	}

	var results []Result
	for _, v := range s.Vars {
		val, ok := env[v.Key]
		if !ok {
			results = append(results, Result{Key: v.Key, Skipped: true, Format: opts.Format})
			continue
		}
		if opts.SensitiveOnly && !v.Sensitive {
			results = append(results, Result{Key: v.Key, Original: val, Encoded: val, Skipped: true, Format: opts.Format})
			continue
		}
		encoded := applyFormat(val, opts.Format)
		results = append(results, Result{
			Key:      v.Key,
			Original: val,
			Encoded:  encoded,
			Format:   opts.Format,
		})
	}
	return results, nil
}

func applyFormat(val string, f Format) string {
	switch f {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(val))
	case FormatHex:
		var sb strings.Builder
		for _, b := range []byte(val) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String()
	default:
		return val
	}
}
