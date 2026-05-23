package fieldwindow

import "encoding/json"

// Rule defines a sliding-window aggregation over a numeric field.
type Rule struct {
	SourceField string
	DestField   string
	Size        int // number of recent values to keep
}

// Transformer maintains per-rule windows and emits a running average.
type Transformer struct {
	rules   []Rule
	windows [][]float64
}

// New returns a Transformer for the given rules.
// Rules with Size < 1 are silently skipped.
func New(rules []Rule) *Transformer {
	var valid []Rule
	var wins [][]float64
	for _, r := range rules {
		if r.Size < 1 {
			continue
		}
		valid = append(valid, r)
		wins = append(wins, make([]float64, 0, r.Size))
	}
	return &Transformer{rules: valid, windows: wins}
}

// Transform adds the source field value to its window and writes the
// rolling average to the destination field. Entries missing the source
// field or with non-numeric values are passed through unchanged for that
// rule.
func (t *Transformer) Transform(entry map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for i, r := range t.rules {
		raw, ok := out[r.SourceField]
		if !ok {
			continue
		}
		var f float64
		if err := json.Unmarshal(raw, &f); err != nil {
			continue
		}
		win := append(t.windows[i], f)
		if len(win) > r.Size {
			win = win[len(win)-r.Size:]
		}
		t.windows[i] = win
		avg := average(win)
		encoded, err := json.Marshal(avg)
		if err != nil {
			continue
		}
		out[r.DestField] = json.RawMessage(encoded)
	}
	return out
}

func average(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
