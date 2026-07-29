# Benchmark results

Measured on 2026-07-28.

## Environment

- CPU: AMD Ryzen 5 5600 6-Core Processor
- OS: Linux
- Architecture: amd64
- Go: 1.26.5
- Connections per benchmark operation: 1,000
- Worker-pool size: 100

Both benchmarks connect to the same temporary TCP listener on
`127.0.0.1`. Reusing one open local port isolates worker-pool coordination
from differences between remote hosts, firewalls, and port states.

## Command

```bash
go test -run '^$' \
  -bench 'Benchmark(Sequential|WorkerPool)Scan$' \
  -benchtime=3x \
  -benchmem \
  -count=3 \
  ./internal/scanner
```

## Raw results

```text
BenchmarkSequentialScan-12    3  52285223 ns/op  1808032 B/op  33001 allocs/op
BenchmarkSequentialScan-12    3  51206250 ns/op  1808037 B/op  33001 allocs/op
BenchmarkSequentialScan-12    3  49143515 ns/op  1809896 B/op  33002 allocs/op
BenchmarkWorkerPoolScan-12    3   7750200 ns/op  1678562 B/op  28662 allocs/op
BenchmarkWorkerPoolScan-12    3   7898278 ns/op  1617496 B/op  28213 allocs/op
BenchmarkWorkerPoolScan-12    3   7650779 ns/op  1623592 B/op  28333 allocs/op
```

## Summary

| Method | Mean time | Mean memory | Mean allocations | Relative speed |
| --- | ---: | ---: | ---: | ---: |
| Sequential | 50.88 ms | 1.81 MB | 33,001 | 1.0× |
| Worker pool | 7.77 ms | 1.64 MB | 28,403 | 6.5× |

## End-to-end localhost scan

The compiled CLI scanned ports `1-1000` on `127.0.0.1` with 100 workers and a
500 ms per-connection timeout:

```bash
scanner \
  -target 127.0.0.1 \
  -ports 1-1000 \
  -workers 100 \
  -timeout 500ms \
  -output json
```

Observed elapsed time:

```text
0.009 seconds
```

Local results are fast because closed loopback ports normally reject
connections immediately. Remote or filtered targets can take much longer and
are bounded by the configured timeout.
