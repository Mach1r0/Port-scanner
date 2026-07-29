package main

import (
	"bytes"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestRunTableOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := run(
		[]string{
			"-target", "127.0.0.1",
			"-ports", availableClosedPort(t),
			"-workers", "1",
			"-timeout", "100ms",
			"-output", "table",
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "PORT") {
		t.Fatalf("table output is missing its header:\n%s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", stderr.String())
	}
}

func TestRunJSONOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := run(
		[]string{
			"-target", "127.0.0.1",
			"-ports", availableClosedPort(t),
			"-workers", "1",
			"-timeout", "100ms",
			"-output", "json",
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
}

func TestRunVerboseProgress(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := run(
		[]string{
			"-target", "127.0.0.1",
			"-ports", availableClosedPort(t),
			"-workers", "1",
			"-timeout", "100ms",
			"-verbose",
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stderr.String(), "[1/1]") {
		t.Fatalf("verbose output is missing progress:\n%s", stderr.String())
	}
}

func TestRunRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "missing target",
			args:     nil,
			contains: "target is required",
		},
		{
			name:     "invalid target",
			args:     []string{"-target", "bad target"},
			contains: "invalid hostname",
		},
		{
			name:     "invalid ports",
			args:     []string{"-target", "127.0.0.1", "-ports", "0"},
			contains: "invalid ports",
		},
		{
			name:     "invalid workers",
			args:     []string{"-target", "127.0.0.1", "-workers", "0"},
			contains: "workers must be greater than zero",
		},
		{
			name:     "invalid timeout",
			args:     []string{"-target", "127.0.0.1", "-timeout", "0s"},
			contains: "timeout must be greater than zero",
		},
		{
			name:     "invalid output",
			args:     []string{"-target", "127.0.0.1", "-output", "yaml"},
			contains: "unsupported output format",
		},
		{
			name:     "unexpected argument",
			args:     []string{"-target", "127.0.0.1", "extra"},
			contains: "unexpected arguments",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			err := run(test.args, &stdout, &stderr)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), test.contains) {
				t.Fatalf("error %q does not contain %q", err, test.contains)
			}
		})
	}
}

func availableClosedPort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("close local listener: %v", err)
	}

	return strconv.Itoa(port)
}
