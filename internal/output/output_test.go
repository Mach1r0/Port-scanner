package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Mach1r0/port-scanner/internal/scanner"
)

func sampleResults() []scanner.Result {
	return []scanner.Result{
		{
			Port:     22,
			Status:   scanner.StatusClosed,
			Duration: 250 * time.Microsecond,
		},
		{
			Port:     80,
			Status:   scanner.StatusOpen,
			Duration: 2 * time.Millisecond,
		},
	}
}

func TestWriteTable(t *testing.T) {
	var buffer bytes.Buffer

	if err := WriteTable(&buffer, sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buffer.String()

	expectedValues := []string{
		"PORT",
		"STATUS",
		"DURATION",
		"22",
		"closed",
		"80",
		"open",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(got, expected) {
			t.Errorf("output does not contain %q:\n%s", expected, got)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	var buffer bytes.Buffer

	if err := WriteJSON(&buffer, sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []jsonResult
	if err := json.Unmarshal(buffer.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d results, want 2", len(got))
	}

	if got[0].Port != 22 {
		t.Errorf("first port is %d, want 22", got[0].Port)
	}

	if got[1].Status != scanner.StatusOpen {
		t.Errorf("second status is %q, want %q", got[1].Status, scanner.StatusOpen)
	}

	if got[1].Duration != "2ms" {
		t.Errorf("second duration is %q, want %q", got[1].Duration, "2ms")
	}
}
