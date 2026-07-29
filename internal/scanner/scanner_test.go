package scanner

import (
	"net"
	"reflect"
	"strings"
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

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{name: "IPv4", target: "127.0.0.1"},
		{name: "IPv6", target: "::1"},
		{name: "localhost", target: "localhost"},
		{name: "hostname", target: "scanner.example.com"},
		{name: "fully qualified hostname", target: "scanner.example.com."},
		{name: "empty", target: "", wantErr: true},
		{name: "whitespace", target: "   ", wantErr: true},
		{name: "space in hostname", target: "bad target", wantErr: true},
		{name: "underscore", target: "bad_target", wantErr: true},
		{name: "leading hyphen", target: "-scanner.example", wantErr: true},
		{name: "empty label", target: "scanner..example", wantErr: true},
		{name: "invalid IPv4", target: "999.999.999.999", wantErr: true},
		{name: "invalid IPv6", target: "2001:db8:::1", wantErr: true},
		{name: "target with port", target: "127.0.0.1:80", wantErr: true},
		{
			name:    "label too long",
			target:  strings.Repeat("a", 64) + ".example",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateTarget(test.target)
			if test.wantErr && err == nil {
				t.Fatal("expected an error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestScanPortOpen(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	port := listener.Addr().(*net.TCPAddr).Port
	result := ScanPort("127.0.0.1", port, time.Second)

	if result.Status != StatusOpen {
		t.Fatalf("got status %q, want %q", result.Status, StatusOpen)
	}
}

func TestScanPortClosed(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on local port: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("close local listener: %v", err)
	}

	result := ScanPort("127.0.0.1", port, time.Second)
	if result.Status != StatusClosed {
		t.Fatalf("got status %q, want %q", result.Status, StatusClosed)
	}
}

func TestScanReportsProgress(t *testing.T) {
	ports := []int{1, 2, 3}
	var calls int

	results, err := ScanWithProgress(
		"127.0.0.1",
		ports,
		2,
		100*time.Millisecond,
		func(completed, total int, _ Result) {
			calls++
			if completed != calls {
				t.Errorf("completed is %d, want %d", completed, calls)
			}
			if total != len(ports) {
				t.Errorf("total is %d, want %d", total, len(ports))
			}
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != len(ports) {
		t.Fatalf("got %d results, want %d", len(results), len(ports))
	}
	if calls != len(ports) {
		t.Fatalf("got %d progress calls, want %d", calls, len(ports))
	}
}

func TestScanRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		ports   []int
		workers int
		timeout time.Duration
	}{
		{
			name:    "invalid target",
			host:    "bad target",
			ports:   []int{80},
			workers: 1,
			timeout: time.Second,
		},
		{
			name:    "invalid port",
			host:    "127.0.0.1",
			ports:   []int{0},
			workers: 1,
			timeout: time.Second,
		},
		{
			name:    "invalid timeout",
			host:    "127.0.0.1",
			ports:   []int{80},
			workers: 1,
			timeout: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Scan(test.host, test.ports, test.workers, test.timeout)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
