package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/levelfilter"
)

func TestParseLevelKnown(t *testing.T) {
	cases := []struct {
		input string
		want  levelfilter.Level
	}{
		{"trace", levelfilter.LevelTrace},
		{"DEBUG", levelfilter.LevelDebug},
		{" Info ", levelfilter.LevelInfo},
		{"WARN", levelfilter.LevelWarn},
		{"error", levelfilter.LevelError},
		{"fatal", levelfilter.LevelFatal},
	}
	for _, c := range cases {
		if got := levelfilter.ParseLevel(c.input); got != c.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestParseLevelUnknown(t *testing.T) {
	if got := levelfilter.ParseLevel("verbose"); got != levelfilter.LevelUnknown {
		t.Errorf("expected LevelUnknown, got %v", got)
	}
}

func TestFilterMatchAboveMin(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn, "level")
	entry := map[string]any{"level": "error", "msg": "boom"}
	if !f.Match(entry) {
		t.Error("expected match for error >= warn")
	}
}

func TestFilterRejectsBelowMin(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn, "level")
	entry := map[string]any{"level": "debug", "msg": "verbose"}
	if f.Match(entry) {
		t.Error("expected no match for debug < warn")
	}
}

func TestFilterPassesMissingField(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	entry := map[string]any{"msg": "no level field"}
	if !f.Match(entry) {
		t.Error("expected pass-through for missing level field")
	}
}

func TestFilterPassesUnknownLevel(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	entry := map[string]any{"level": "verbose"}
	if !f.Match(entry) {
		t.Error("expected pass-through for unrecognised level value")
	}
}

func TestFilterDefaultField(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelInfo, "")
	entry := map[string]any{"level": "debug"}
	if f.Match(entry) {
		t.Error("expected rejection using default 'level' field")
	}
}

func TestFilterNonStringLevel(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	entry := map[string]any{"level": 42}
	if !f.Match(entry) {
		t.Error("expected pass-through for non-string level value")
	}
}
