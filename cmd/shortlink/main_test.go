package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"version"}, &out, &errOut); code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errOut.String())
	}
	if !strings.HasPrefix(out.String(), "shortlink ") {
		t.Fatalf("unexpected output %q", out.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"nope"}, &out, &errOut); code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
}

func TestRunServeBadData(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "links.json")
	if err := os.WriteFile(bad, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		data string
	}{
		{"directory", dir},
		{"unparseable", bad},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := run([]string{"serve", "--addr", "127.0.0.1:0", "--data", tt.data}, &out, &errOut)
			if code != 1 {
				t.Fatalf("exit %d, want 1", code)
			}
			if !strings.Contains(errOut.String(), tt.data) {
				t.Fatalf("stderr %q does not name %q", errOut.String(), tt.data)
			}
		})
	}
}
