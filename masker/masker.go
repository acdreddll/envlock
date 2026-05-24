package masker

import (
	"strings"

	"github.com/envlock/schema"
)

// Options controls masking behaviour.
type Options struct {
	// SensitiveOnly restricts masking to vars marked sensitive=true.
	SensitiveOnly bool
	// MaskChar is the character used to fill the masked portion (default '*').
	MaskChar rune
	// VisiblePrefix is the number of leading characters to leave visible.
	VisiblePrefix int
	// VisibleSuffix is the number of trailing characters to leave visible.
	VisibleSuffix int
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		SensitiveOnly: true,
		MaskChar:      '*',
		VisiblePrefix: 0,
		VisibleSuffix: 0,
	}
}

// Result holds a single variable's masked output.
type Result struct {
	Key      string
	Masked   string
	WasMasked bool
}

// Mask applies masking rules to the provided env map using the schema.
// It returns one Result per schema variable that appears in env.
func Mask(s schema.Schema, env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(s.Vars))
	for _, v := range s.Vars {
		val, ok := env[v.Key]
		if !ok {
			continue
		}
		shouldMask := !opts.SensitiveOnly || v.Sensitive
		if shouldMask {
			results = append(results, Result{
				Key:       v.Key,
				Masked:    applyMask(val, opts),
				WasMasked: true,
			})
		} else {
			results = append(results, Result{
				Key:       v.Key,
				Masked:    val,
				WasMasked: false,
			})
		}
	}
	return results
}

func applyMask(val string, opts Options) string {
	if len(val) == 0 {
		return val
	}
	pre := opts.VisiblePrefix
	suf := opts.VisibleSuffix
	if pre+suf >= len(val) {
		// Mask everything when the window is too wide.
		return strings.Repeat(string(opts.MaskChar), len(val))
	}
	midLen := len(val) - pre - suf
	return val[:pre] + strings.Repeat(string(opts.MaskChar), midLen) + val[len(val)-suf:]
}
