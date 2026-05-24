package fieldgeo_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldgeo"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func rules() []fieldgeo.Rule {
	return []fieldgeo.Rule{
		{
			SrcField:  "ip",
			DestField: "region",
			Regions: map[string]string{
				"10.0.0.0/8":     "private",
				"192.168.0.0/16": "lan",
				"8.8.8.0/24":     "google-dns",
			},
		},
	}
}

func TestIdentityNoRules(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(nil)
	in := entry("ip", "10.1.2.3")
	out := tr.Apply(in)
	if _, ok := out["region"]; ok {
		t.Fatal("expected no region field")
	}
}

func TestAnnotatesMatchingCIDR(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	out := tr.Apply(entry("ip", "10.1.2.3"))
	if out["region"] != "private" {
		t.Fatalf("expected private, got %v", out["region"])
	}
}

func TestAnnotatesLAN(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	out := tr.Apply(entry("ip", "192.168.1.55"))
	if out["region"] != "lan" {
		t.Fatalf("expected lan, got %v", out["region"])
	}
}

func TestNoMatchLeavesEntryUnchanged(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	out := tr.Apply(entry("ip", "1.2.3.4"))
	if _, ok := out["region"]; ok {
		t.Fatal("expected no region field for unmatched IP")
	}
}

func TestMissingSrcFieldSkipped(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	out := tr.Apply(entry("msg", "hello"))
	if _, ok := out["region"]; ok {
		t.Fatal("expected no region field when src field missing")
	}
}

func TestInvalidIPSkipped(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	out := tr.Apply(entry("ip", "not-an-ip"))
	if _, ok := out["region"]; ok {
		t.Fatal("expected no region field for invalid IP")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	t.Parallel()
	tr := fieldgeo.New(rules())
	in := entry("ip", "10.0.0.1")
	_ = tr.Apply(in)
	if _, ok := in["region"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}
