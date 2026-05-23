package classifier

import (
	"strings"

	"github.com/envlock/schema"
)

// Category represents a classification bucket for an environment variable.
type Category string

const (
	CategorySecret      Category = "secret"
	CategoryConfig      Category = "config"
	CategoryFeatureFlag Category = "feature_flag"
	CategoryInfra       Category = "infra"
	CategoryUnknown     Category = "unknown"
)

// Result holds the classification of a single variable.
type Result struct {
	Key      string
	Category Category
	Reason   string
}

// Classify assigns a Category to each variable in the schema based on its
// metadata: sensitive flag, key prefix/suffix patterns, and tags.
func Classify(s schema.Schema) []Result {
	results := make([]Result, 0, len(s.Vars))
	for _, v := range s.Vars {
		results = append(results, classifyVar(v))
	}
	return results
}

func classifyVar(v schema.EnvVar) Result {
	key := strings.ToUpper(v.Key)

	if v.Sensitive {
		return Result{Key: v.Key, Category: CategorySecret, Reason: "marked sensitive"}
	}

	for _, tag := range v.Tags {
		switch strings.ToLower(tag) {
		case "secret", "credential", "auth":
			return Result{Key: v.Key, Category: CategorySecret, Reason: "tag: " + tag}
		case "feature", "flag", "toggle":
			return Result{Key: v.Key, Category: CategoryFeatureFlag, Reason: "tag: " + tag}
		case "infra", "infrastructure", "platform":
			return Result{Key: v.Key, Category: CategoryInfra, Reason: "tag: " + tag}
		}
	}

	switch {
	case hasAnyPrefix(key, "SECRET_", "TOKEN_", "API_KEY_", "PASSWORD_", "PRIVATE_"):
		return Result{Key: v.Key, Category: CategorySecret, Reason: "key prefix"}
	case hasAnySuffix(key, "_SECRET", "_TOKEN", "_PASSWORD", "_KEY", "_CERT"):
		return Result{Key: v.Key, Category: CategorySecret, Reason: "key suffix"}
	case hasAnyPrefix(key, "ENABLE_", "DISABLE_", "FEATURE_", "FLAG_"):
		return Result{Key: v.Key, Category: CategoryFeatureFlag, Reason: "key prefix"}
	case hasAnyPrefix(key, "HOST_", "PORT_", "ADDR_", "ENDPOINT_", "REGION_", "CLUSTER_"):
		return Result{Key: v.Key, Category: CategoryInfra, Reason: "key prefix"}
	case hasAnySuffix(key, "_HOST", "_PORT", "_ADDR", "_ENDPOINT", "_URL", "_URI"):
		return Result{Key: v.Key, Category: CategoryInfra, Reason: "key suffix"}
	default:
		return Result{Key: v.Key, Category: CategoryConfig, Reason: "default"}
	}
}

func hasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, sfx := range suffixes {
		if strings.HasSuffix(s, sfx) {
			return true
		}
	}
	return false
}
