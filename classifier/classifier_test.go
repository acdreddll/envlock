package classifier_test

import (
	"testing"

	"github.com/envlock/classifier"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string) schema.EnvVar { return schema.EnvVar{Key: key} }

func TestClassify_SensitiveFlag(t *testing.T) {
	s := makeSchema(schema.EnvVar{Key: "DB_PASS", Sensitive: true})
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategorySecret {
		t.Errorf("expected secret, got %s", res[0].Category)
	}
	if res[0].Reason != "marked sensitive" {
		t.Errorf("unexpected reason: %s", res[0].Reason)
	}
}

func TestClassify_SecretByPrefix(t *testing.T) {
	s := makeSchema(ev("TOKEN_GITHUB"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategorySecret {
		t.Errorf("expected secret, got %s", res[0].Category)
	}
}

func TestClassify_SecretBySuffix(t *testing.T) {
	s := makeSchema(ev("STRIPE_KEY"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategorySecret {
		t.Errorf("expected secret, got %s", res[0].Category)
	}
}

func TestClassify_FeatureFlag(t *testing.T) {
	s := makeSchema(ev("ENABLE_DARK_MODE"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategoryFeatureFlag {
		t.Errorf("expected feature_flag, got %s", res[0].Category)
	}
}

func TestClassify_InfraByPrefix(t *testing.T) {
	s := makeSchema(ev("HOST_REDIS"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategoryInfra {
		t.Errorf("expected infra, got %s", res[0].Category)
	}
}

func TestClassify_InfraBySuffix(t *testing.T) {
	s := makeSchema(ev("DATABASE_URL"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategoryInfra {
		t.Errorf("expected infra, got %s", res[0].Category)
	}
}

func TestClassify_DefaultConfig(t *testing.T) {
	s := makeSchema(ev("LOG_LEVEL"))
	res := classifier.Classify(s)
	if res[0].Category != classifier.CategoryConfig {
		t.Errorf("expected config, got %s", res[0].Category)
	}
}

func TestClassify_TagOverridesPrefix(t *testing.T) {
	v := schema.EnvVar{Key: "HOST_SOMETHING", Tags: []string{"feature"}}
	res := classifier.Classify(makeSchema(v))
	if res[0].Category != classifier.CategoryFeatureFlag {
		t.Errorf("expected feature_flag from tag, got %s", res[0].Category)
	}
}

func TestClassify_EmptySchema(t *testing.T) {
	res := classifier.Classify(makeSchema())
	if len(res) != 0 {
		t.Errorf("expected empty results, got %d", len(res))
	}
}
