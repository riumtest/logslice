// Package fieldhasher computes a hash of one or more source fields and
// writes the result into a destination field. This is useful for generating
// stable identifiers or anonymising values before forwarding logs.
package fieldhasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

// Algorithm selects the hash function used.
type Algorithm string

const (
	MD5    Algorithm = "md5"
	SHA1   Algorithm = "sha1"
	SHA256 Algorithm = "sha256"
)

// Rule describes a single hashing operation.
type Rule struct {
	// Fields are the source fields whose string values are concatenated before hashing.
	Fields []string
	// Dest is the field written with the hex-encoded digest.
	Dest string
	// Algo selects the hash algorithm; defaults to SHA256.
	Algo Algorithm
}

// Hasher applies hashing rules to log entries.
type Hasher struct {
	rules []Rule
}

// New returns a Hasher that applies the given rules in order.
func New(rules []Rule) *Hasher {
	return &Hasher{rules: rules}
}

// Apply processes a single log entry and returns the (possibly modified) copy.
func (h *Hasher) Apply(entry map[string]any) map[string]any {
	if len(h.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range h.rules {
		if r.Dest == "" || len(r.Fields) == 0 {
			continue
		}
		hw := newHash(r.Algo)
		for _, f := range r.Fields {
			v, ok := out[f]
			if !ok {
				continue
			}
			fmt.Fprintf(hw, "%v", v)
		}
		out[r.Dest] = hex.EncodeToString(hw.Sum(nil))
	}
	return out
}

func newHash(algo Algorithm) hash.Hash {
	switch algo {
	case MD5:
		return md5.New()
	case SHA1:
		return sha1.New()
	default:
		return sha256.New()
	}
}

func shallowCopy(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
