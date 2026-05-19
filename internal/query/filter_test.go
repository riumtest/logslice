package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAndMatch(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		entry   map[string]any
		want    bool
	}{
		{"eq match", "level=error", map[string]any{"level": "error"}, true},
		{"eq no match", "level=error", map[string]any{"level": "info"}, false},
		{"neq match", "level!=info", map[string]any{"level": "error"}, true},
		{"contains match", "msg~timeout", map[string]any{"msg": "connection timeout occurred"}, true},
		{"contains no match", "msg~timeout", map[string]any{"msg": "all good"}, false},
		{"gt match", "status>399", map[string]any{"status": float64(500)}, true},
		{"gt no match", "status>399", map[string]any{"status": float64(200)}, false},
		{"gte match", "status>=400", map[string]any{"status": float64(400)}, true},
		{"lt match", "latency<100", map[string]any{"latency": float64(50)}, true},
		{"lte match", "latency<=100", map[string]any{"latency": float64(100)}, true},
		{"missing field", "level=error", map[string]any{"msg": "hi"}, false},
		{"multi filter all match", "level=error status>399", map[string]any{"level": "error", "status": float64(500)}, true},
		{"multi filter partial", "level=error status>399", map[string]any{"level": "error", "status": float64(200)}, false},
		{"empty query matches all", "", map[string]any{"level": "info"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := Parse(tt.query)
			require.NoError(t, err)
			assert.Equal(t, tt.want, q.Match(tt.entry))
		})
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse("badtoken")
	assert.Error(t, err)
}

func TestParseEmptyFieldOrValue(t *testing.T) {
	_, err := parseFilter("=value")
	assert.Error(t, err)

	_, err = parseFilter("field=")
	assert.Error(t, err)
}

func TestMatchEmptyEntry(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"eq against empty entry", "level=error"},
		{"numeric gt against empty entry", "status>399"},
		{"contains against empty entry", "msg~timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := Parse(tt.query)
			require.NoError(t, err)
			assert.False(t, q.Match(map[string]any{}))
		})
	}
}
