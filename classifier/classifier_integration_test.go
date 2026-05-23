package classifier_test

import (
	"testing"

	"github.com/envlock/classifier"
	"github.com/envlock/schema"
)

func TestClassify_MixedSchema(t *testing.T) {
	s := schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "APP_NAME"},
			{Key: "DATABASE_URL"},
			{Key: "STRIPE_SECRET_KEY", Sensitive: true},
			{Key: "ENABLE_BETA_FEATURES"},
			{Key: "REDIS_HOST"},
			{Key: "JWT_TOKEN"},
		},
	}

	results := classifier.Classify(s)

	expected := map[string]classifier.Category{
		"APP_NAME":             classifier.CategoryConfig,
		"DATABASE_URL":         classifier.CategoryInfra,
		"STRIPE_SECRET_KEY":    classifier.CategorySecret,
		"ENABLE_BETA_FEATURES": classifier.CategoryFeatureFlag,
		"REDIS_HOST":           classifier.CategoryInfra,
		"JWT_TOKEN":            classifier.CategorySecret,
	}

	if len(results) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(results))
	}

	for _, r := range results {
		want, ok := expected[r.Key]
		if !ok {
			t.Errorf("unexpected key in results: %s", r.Key)
			continue
		}
		if r.Category != want {
			t.Errorf("key %s: expected %s, got %s (reason: %s)", r.Key, want, r.Category, r.Reason)
		}
	}
}
