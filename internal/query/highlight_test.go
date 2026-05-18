package query

import (
	"testing"
)

func TestHighlighterNoTerms(t *testing.T) {
	h := NewHighlighter(nil)
	if h.HasTerms() {
		t.Fatal("expected no terms")
	}
	got := h.Wrap("hello world", "[", "]")
	if got != "hello world" {
		t.Fatalf("expected unchanged string, got %q", got)
	}
}

func TestHighlighterEmptyTermsFiltered(t *testing.T) {
	h := NewHighlighter([]string{"", ""})
	if h.HasTerms() {
		t.Fatal("expected no terms after filtering empties")
	}
}

func TestHighlighterWrapCaseInsensitive(t *testing.T) {
	h := NewHighlighter([]string{"error"})
	input := "An ERROR occurred in the service"
	got := h.Wrap(input, "[", "]")
	want := "An [error] occurred in the service"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestHighlighterMultipleTerms(t *testing.T) {
	h := NewHighlighter([]string{"warn", "timeout"})
	input := "WARN: connection timeout"
	got := h.Wrap(input, "<", ">")
	want := "<warn>: connection <timeout>"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestHighlighterNoMatch(t *testing.T) {
	h := NewHighlighter([]string{"fatal"})
	input := "everything is fine"
	got := h.Wrap(input, "[", "]")
	if got != input {
		t.Fatalf("expected unchanged string, got %q", got)
	}
}

func TestReplaceInsensitiveEmptyOld(t *testing.T) {
	got := replaceInsensitive("hello", "", "X")
	if got != "hello" {
		t.Fatalf("expected unchanged string, got %q", got)
	}
}
