package templater

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/yourorg/envlock/schema"
)

// Result holds the rendered template output.
type Result struct {
	Output string
	Missing []string
}

// Render takes a Go template string and substitutes {{.ENV_VAR}} placeholders
// using values from the provided environment map, validated against the schema.
func Render(tmpl string, env map[string]string, s *schema.Schema) (*Result, error) {
	result := &Result{}

	// Build substitution map: use env value or default, track missing required vars.
	subst := make(map[string]string)
	for _, entry := range s.Vars {
		val, ok := env[entry.Key]
		if !ok || val == "" {
			if entry.Default != "" {
				subst[entry.Key] = entry.Default
			} else if entry.Required {
				result.Missing = append(result.Missing, entry.Key)
				subst[entry.Key] = ""
			} else {
				subst[entry.Key] = ""
			}
		} else {
			subst[entry.Key] = val
		}
	}

	// Also pass through any env keys not in schema.
	for k, v := range env {
		if _, exists := subst[k]; !exists {
			subst[k] = v
		}
	}

	// Convert map keys to template-safe identifiers.
	funcMap := template.FuncMap{
		"env": func(key string) string {
			return subst[key]
		},
	}

	// Replace {{KEY}} style with {{env "KEY"}} for convenience.
	normalized := normalizePlaceholders(tmpl)

	t, err := template.New("envlock").Funcs(funcMap).Parse(normalized)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, subst); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}

	result.Output = buf.String()
	return result, nil
}

// normalizePlaceholders converts ${VAR} style placeholders to {{env "VAR"}}.
func normalizePlaceholders(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			end := strings.Index(s[i:], "}")
			if end != -1 {
				key := s[i+2 : i+end]
				b.WriteString(`{{env "`)
				b.WriteString(key)
				b.WriteString(`"}}`)
				i += end + 1
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
