// Package fieldgeo annotates log entries with geographic region labels
// derived from a simple IP-to-region lookup table supplied by the caller.
package fieldgeo

import (
	"net"
)

// Rule describes a single IP-to-region annotation.
type Rule struct {
	// SrcField is the entry field that contains an IP address string.
	SrcField string
	// DestField is written with the resolved region string.
	DestField string
	// Regions maps CIDR prefix strings to region labels.
	Regions map[string]string
}

// compiledRule holds a parsed version of a Rule.
type compiledRule struct {
	Rule
	nets []*net.IPNet
	labels []string
}

// Transformer applies geographic annotation rules to log entries.
type Transformer struct {
	rules []compiledRule
}

// New builds a Transformer from the provided rules.
// CIDR strings that cannot be parsed are silently skipped.
func New(rules []Rule) *Transformer {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		cr := compiledRule{Rule: r}
		for cidr, label := range r.Regions {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			cr.nets = append(cr.nets, ipNet)
			cr.labels = append(cr.labels, label)
		}
		compiled = append(compiled, cr)
	}
	return &Transformer{rules: compiled}
}

// Apply annotates entry in-place and returns it.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	for _, cr := range t.rules {
		raw, ok := entry[cr.SrcField]
		if !ok {
			continue
		}
		ipStr, ok := raw.(string)
		if !ok {
			continue
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}
		for i, ipNet := range cr.nets {
			if ipNet.Contains(ip) {
				out := shallowCopy(entry)
				out[cr.DestField] = cr.labels[i]
				entry = out
				break
			}
		}
	}
	return entry
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
