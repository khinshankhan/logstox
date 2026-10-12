package zapx

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/khinshankhan/logstox"
)

// captureStderr builds a logger with b and logs one line through it. The
// default sinks are stderr, so it's swapped for a pipe while the logger is
// built and used.
func captureStderr(t *testing.T, b Backend) map[string]any {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stderr
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = orig
	})

	lg := b.New(logstox.Options[ZapField]{})
	lg.Info("hello")
	_ = lg.Sync()
	w.Close()
	os.Stderr = orig

	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		t.Fatalf("read log line: %v", err)
	}
	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("decode %q: %v", line, err)
	}
	return entry
}

func TestCaller(t *testing.T) {
	tests := []struct {
		name       string
		addSource  bool
		wantCaller bool
	}{
		{"off unless asked", false, false},
		{"points at the call site when asked", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// production config, so the line is JSON
			entry := captureStderr(t, Backend{AddSource: tt.addSource})

			caller, ok := entry["caller"].(string)
			if ok != tt.wantCaller {
				t.Fatalf("caller = %q (present=%v), want present=%v", caller, ok, tt.wantCaller)
			}
			// the line in captureStderr that logged, not this package's wrapper
			if ok && !strings.HasPrefix(caller, "zapx/zapx_test.go:") {
				t.Errorf("caller = %q, want the call site in zapx_test.go", caller)
			}
		})
	}
}
