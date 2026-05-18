package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/formatter"
)

func TestFormatterText(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf)

	entry := map[string]interface{}{
		"time":  "2024-01-15T10:30:00Z",
		"level": "info",
		"msg":   "server started",
		"port":  8080,
	}

	if err := f.Write(entry, `{"raw":true}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", out)
	}
	if !strings.Contains(out, "server started") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port field in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-15 10:30:00") {
		t.Errorf("expected formatted time in output, got: %s", out)
	}
}

func TestFormatterJSON(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.WithFormat(formatter.FormatJSON))

	raw := `{"level":"error","msg":"oops"}`
	if err := f.Write(nil, raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.TrimSpace(buf.String()) != raw {
		t.Errorf("expected raw JSON output, got: %s", buf.String())
	}
}

func TestFormatterCustomKeys(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.WithKeys("timestamp", "severity", "message"))

	entry := map[string]interface{}{
		"timestamp": "2024-06-01T00:00:00Z",
		"severity":  "warn",
		"message":   "disk usage high",
	}

	if err := f.Write(entry, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
	if !strings.Contains(out, "disk usage high") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestFormatterMissingFields(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf)

	entry := map[string]interface{}{
		"service": "auth",
		"user_id": 42,
	}

	if err := f.Write(entry, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "service=auth") {
		t.Errorf("expected service field, got: %s", out)
	}
}
