package hasher_test

import (
	"testing"

	"github.com/envlock/hasher"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, opts ...func(*schema.EnvVar)) schema.EnvVar {
	v := schema.EnvVar{Key: key}
	for _, o := range opts {
		o(&v)
	}
	return v
}

func withDesc(d string) func(*schema.EnvVar) {
	return func(v *schema.EnvVar) { v.Description = d }
}

func required(v *schema.EnvVar)  { v.Required = true }
func sensitive(v *schema.EnvVar) { v.Sensitive = true }

func TestHash_EmptySchema(t *testing.T) {
	s := makeSchema()
	sum := hasher.Hash(s)
	if sum.SchemaHash == "" {
		t.Fatal("expected non-empty schema hash for empty schema")
	}
	if len(sum.Vars) != 0 {
		t.Fatalf("expected 0 var results, got %d", len(sum.Vars))
	}
}

func TestHash_SingleVar(t *testing.T) {
	s := makeSchema(ev("PORT", withDesc("HTTP port"), required))
	sum := hasher.Hash(s)
	if len(sum.Vars) != 1 {
		t.Fatalf("expected 1 var result, got %d", len(sum.Vars))
	}
	if sum.Vars[0].Key != "PORT" {
		t.Errorf("unexpected key: %s", sum.Vars[0].Key)
	}
	if len(sum.Vars[0].Hash) != 64 {
		t.Errorf("expected 64-char hex hash, got len %d", len(sum.Vars[0].Hash))
	}
}

func TestHash_DeterministicAcrossCallOrder(t *testing.T) {
	// Schema with vars in different insertion order should yield same schema hash.
	s1 := makeSchema(ev("A"), ev("B"), ev("C"))
	s2 := makeSchema(ev("C"), ev("A"), ev("B"))

	sum1 := hasher.Hash(s1)
	sum2 := hasher.Hash(s2)

	if sum1.SchemaHash != sum2.SchemaHash {
		t.Errorf("schema hashes differ across insertion orders: %s vs %s",
			sum1.SchemaHash, sum2.SchemaHash)
	}
}

func TestHash_ChangeDetected(t *testing.T) {
	s1 := makeSchema(ev("SECRET", sensitive, withDesc("original")))
	s2 := makeSchema(ev("SECRET", sensitive, withDesc("changed")))

	sum1 := hasher.Hash(s1)
	sum2 := hasher.Hash(s2)

	if sum1.SchemaHash == sum2.SchemaHash {
		t.Error("expected different schema hash after description change")
	}
	if sum1.Vars[0].Hash == sum2.Vars[0].Hash {
		t.Error("expected different var hash after description change")
	}
}

func TestHash_VarsAreSortedByKey(t *testing.T) {
	s := makeSchema(ev("ZEBRA"), ev("ALPHA"), ev("MANGO"))
	sum := hasher.Hash(s)
	keys := []string{sum.Vars[0].Key, sum.Vars[1].Key, sum.Vars[2].Key}
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("vars not sorted by key: %v", keys)
	}
}
