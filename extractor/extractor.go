package extractor

import (
	"fmt"
	"regexp"

	"github.com/envlock/schema"
)

// Result holds the extracted subset of schema variables.
type Result struct {
	Vars []schema.EnvVar
}

// Options controls which variables are extracted.
type Options struct {
	// Tags filters by one or more tags (OR logic).
	Tags []string
	// Group filters by group name.
	Group string
	// Pattern filters by key name regexp.
	Pattern string
	// SensitiveOnly restricts to sensitive vars.
	SensitiveOnly bool
	// RequiredOnly restricts to required vars.
	RequiredOnly bool
}

// Extract returns schema variables that match all non-zero criteria in opts.
func Extract(s schema.Schema, opts Options) (Result, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	tagSet := make(map[string]struct{}, len(opts.Tags))
	for _, t := range opts.Tags {
		tagSet[t] = struct{}{}
	}

	var matched []schema.EnvVar
	for _, v := range s.Vars {
		if opts.SensitiveOnly && !v.Sensitive {
			continue
		}
		if opts.RequiredOnly && !v.Required {
			continue
		}
		if opts.Group != "" && v.Group != opts.Group {
			continue
		}
		if re != nil && !re.MatchString(v.Key) {
			continue
		}
		if len(tagSet) > 0 && !hasTag(v.Tags, tagSet) {
			continue
		}
		matched = append(matched, v)
	}
	return Result{Vars: matched}, nil
}

func hasTag(varTags []string, want map[string]struct{}) bool {
	for _, t := range varTags {
		if _, ok := want[t]; ok {
			return true
		}
	}
	return false
}
