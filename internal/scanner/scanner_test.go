package scanner

import (
	"reflect"
	"testing"
	"time"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:  "single port",
			input: "80",
			want:  []int{80},
		},
		{
			name:  "comma separated ports",
			input: "22,80,443",
			want:  []int{22, 80, 443},
		},
		{
			name:  "port range",
			input: "80-83",
			want:  []int{80, 81, 82, 83},
		},
		{
			name:  "combined values",
			input: "22,80-82,443",
			want:  []int{22, 80, 81, 82, 443},
		},
		{
			name:  "duplicates are removed",
			input: "22,20-22,22",
			want:  []int{20, 21, 22},
		},
		{
			name:    "port below minimum",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "port above maximum",
			input:   "65536",
			wantErr: true,
		},
		{
			name:    "reversed range",
			input:   "100-80",
			wantErr: true,
		},
		{
			name:    "invalid text",
			input:   "http",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParsePorts(test.input)

			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestScanReturnsSortedResults(t *testing.T) {
	ports := []int{3, 1, 2}

	results, err := Scan(
		"127.0.0.1",
		ports,
		2,
		100*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(ports) {
		t.Fatalf("got %d results, want %d", len(results), len(ports))
	}

	for i, expectedPort := range []int{1, 2, 3} {
		if results[i].Port != expectedPort {
			t.Fatalf(
				"result %d has port %d, want %d",
				i,
				results[i].Port,
				expectedPort,
			)
		}
	}
}

func TestScanRejectsInvalidWorkerCount(t *testing.T) {
	_, err := Scan(
		"127.0.0.1",
		[]int{80},
		0,
		100*time.Millisecond,
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestScanWithNoPorts(t *testing.T) {
	results, err := Scan(
		"127.0.0.1",
		nil,
		10,
		100*time.Millisecond,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("got %d results, want 0", len(results))
	}
}
