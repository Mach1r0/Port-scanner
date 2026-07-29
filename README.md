# Port Scanner

A concurrent TCP port scanner written in Go. It uses a bounded worker pool,
channels, and `sync.WaitGroup` to scan many ports without creating one
goroutine per port.

> Only scan systems you own or have explicit permission to test.

## Features

- Single ports, comma-separated lists, and ranges
- Configurable worker count and connection timeout
- TCP status classification: `open`, `closed`, or `filtered`
- Deterministic results sorted by port
- Human-readable table and JSON output
- Optional real-time progress on stderr
- Race-tested worker pool
- Multi-stage Docker image and Docker Compose service
- GitHub Actions checks for tests, `go vet`, and `golangci-lint`

## Requirements

- Go 1.26 or newer
- Docker and Docker Compose are optional

## Build

```bash
git clone https://github.com/Mach1r0/port-scanner.git
cd port-scanner
go build -o scanner ./cmd/scanner
```

## Usage

```bash
./scanner -target 127.0.0.1 -ports 1-1000 -workers 100
```

You can also run the source directly:

```bash
go run ./cmd/scanner \
  -target 127.0.0.1 \
  -ports 22,80,443,8000-8010 \
  -workers 20 \
  -timeout 500ms \
  -output table \
  -verbose
```

### Flags

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `-target` | string | required | IPv4, IPv6, or hostname to scan |
| `-ports` | string | `1-1024` | Port, list, range, or combination |
| `-workers` | int | `100` | Maximum concurrent scanning workers |
| `-timeout` | duration | `1s` | Timeout for each TCP connection |
| `-output` | string | `table` | Output format: `table` or `json` |
| `-verbose` | bool | `false` | Print real-time progress to stderr |

Accepted port expressions include:

```text
80
22,80,443
1-1000
22,80,8000-8010
```

Duplicate ports are removed and results are sorted in ascending order.

## Output

### Table

```text
PORT  STATUS  DURATION
22    closed  59.801µs
80    open    1.204ms
443   closed  25.001µs
```

### JSON

```json
[
  {
    "port": 22,
    "status": "closed",
    "duration": "31.92µs"
  },
  {
    "port": 80,
    "status": "open",
    "duration": "1.204ms"
  }
]
```

Verbose progress is written to stderr, so stdout remains valid table or JSON:

```text
Scanning 127.0.0.1: 3 ports with 3 workers and timeout 500ms
[1/3] port 80: open
[2/3] port 22: closed
[3/3] port 443: closed
```

## Status classification

- `open`: the TCP connection completed successfully.
- `filtered`: the connection timed out.
- `closed`: the connection failed without timing out, typically because it was
  refused.

## Architecture

```text
Port producer
     |
     v
Jobs channel ---> N workers ---> Results channel ---> Aggregator ---> Sorted output
                      |
                      v
                  ScanPort
```

The producer closes the jobs channel after sending every port. Workers call
`ScanPort`, send results to a dedicated channel, and signal completion through
`sync.WaitGroup`. A closing goroutine waits for all workers and then closes the
results channel. The aggregator receives every result and sorts the final
slice.

## Docker

Build and run directly:

```bash
docker build -t port-scanner .
docker run --rm --network host port-scanner \
  -target 127.0.0.1 \
  -ports 1-1000 \
  -workers 100 \
  -timeout 500ms
```

On Linux, host networking allows `127.0.0.1` inside the container to refer to
the host. Other Docker platforms may require `host.docker.internal` instead.

Docker Compose has a default localhost scan:

```bash
docker compose up --build
```

Override its arguments with:

```bash
docker compose run --rm scanner \
  -target 127.0.0.1 \
  -ports 22,80,443 \
  -output json
```

## Tests and quality checks

```bash
gofmt -w .
go test -race -cover ./...
go vet ./...
golangci-lint run ./...
```

The tests cover port parsing, validation, open and closed TCP ports, worker
coordination, result sorting, progress reporting, CLI validation, and both
output formats.

## Benchmark

Benchmarks perform 1,000 TCP connections against a temporary local listener.
Results from an AMD Ryzen 5 5600 running Linux:

| Method | Workers | Mean time | Relative speed |
| --- | ---: | ---: | ---: |
| Sequential | 1 | 50.88 ms | 1.0× |
| Worker pool | 100 | 7.77 ms | 6.5× |

A separate scan of localhost ports `1-1000` with 100 workers completed in
`0.009s` on the same machine.

Reproduce the benchmark:

```bash
go test -run '^$' \
  -bench 'Benchmark(Sequential|WorkerPool)Scan$' \
  -benchtime=3x \
  -benchmem \
  -count=3 \
  ./internal/scanner
```

Full measurements and methodology are in
[`benchmark/results.md`](benchmark/results.md).
