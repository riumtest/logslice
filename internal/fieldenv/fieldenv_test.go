package fieldenv_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldenv"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func fakeEnv(env map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldenv.New(nil)
	in := entry("msg", "hello")
	out := tr.Transform(in)
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestInjectsEnvValue(t *testing.T) {
	env := fakeEnv(map[string]string{"APP_ENV": "production"})
	rules := []fieldenv.Rule{{Env: "APP_ENV", Dest: "environment"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(env))
	out := tr.Transform(entry("msg", "hi"))
	if out["environment"] != "production" {
		t.Fatalf("expected environment=production, got %v", out["environment"])
	}
}

func TestUsesDefaultWhenUnset(t *testing.T) {
	env := fakeEnv(map[string]string{})
	rules := []fieldenv.Rule{{Env: "REGION", Dest: "region", Default: "us-east-1"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(env))
	out := tr.Transform(entry("msg", "hi"))
	if out["region"] != "us-east-1" {
		t.Fatalf("expected region=us-east-1, got %v", out["region"])
	}
}

func TestDoesNotOverwriteExistingField(t *testing.T) {
	env := fakeEnv(map[string]string{"APP_ENV": "staging"})
	rules := []fieldenv.Rule{{Env: "APP_ENV", Dest: "env"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(env))
	out := tr.Transform(entry("env", "original"))
	if out["env"] != "original" {
		t.Fatalf("expected env=original, got %v", out["env"])
	}
}

func TestEmptyEnvAndNoDefaultOmitsField(t *testing.T) {
	env := fakeEnv(map[string]string{})
	rules := []fieldenv.Rule{{Env: "MISSING", Dest: "dest"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(env))
	out := tr.Transform(entry("msg", "hi"))
	if _, ok := out["dest"]; ok {
		t.Fatal("expected dest field to be absent")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	env := fakeEnv(map[string]string{"X": "1"})
	rules := []fieldenv.Rule{{Env: "X", Dest: "x"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(env))
	in := entry("msg", "hi")
	tr.Transform(in)
	if _, ok := in["x"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}
