# Port Scanner TODO

This checklist reflects the current implementation and maps back to the
requirements in `README.md`.

## Status legend

- [x] Implemented and verified
- [ ] Missing or incomplete

## Core scanning

- [x] Define scan statuses: `open`, `closed`, and `filtered`
- [x] Define a result containing the port, status, and response duration
- [x] Scan one TCP port with `net.DialTimeout`
- [x] Classify successful connections as open
- [x] Classify connection timeouts as filtered
- [x] Classify other connection errors as closed
- [ ] Validate the target hostname or IP address
- [x] Validate that ports are between 1 and 65535
- [x] Parse a single port, such as `80`
- [x] Parse comma-separated ports, such as `22,80,443`
- [x] Parse port ranges, such as `1-1000`
- [x] Remove duplicate ports and return them in a predictable order

## Worker pool

- [x] Define a configurable scanner or scan function
- [x] Create a jobs channel for ports
- [x] Start a configurable number of workers
- [x] Make each worker call `ScanPort`
- [x] Send worker results through a dedicated results channel
- [x] Coordinate worker completion with `sync.WaitGroup`
- [x] Close channels safely
- [x] Aggregate every scan result
- [x] Sort results by port before returning them
- [x] Avoid starting one goroutine per port

## Command-line interface

The CLI now accepts a target, port expression, worker count, timeout, and
verbose mode.

- [x] Add required `-target` flag
- [x] Add `-ports` flag with a default such as `1-1024`
- [x] Add `-workers` flag with a default such as `100`
- [x] Add `-timeout` flag with a default such as `1s`
- [x] Add `-output` flag accepting `table` or `json`
- [x] Add `-verbose` flag
- [ ] Reject missing or invalid flag values with helpful messages
- [x] Return a non-zero exit code when configuration or scanning fails
- [x] Connect the CLI to the worker-pool scanner

## Output

The output package supports terminal tables and formatted JSON.

- [x] Implement terminal table output
- [x] Implement JSON output
- [x] Include port, status, and response duration in both formats
- [x] Send normal output to stdout
- [x] Send errors and verbose progress to stderr
- [x] Keep silent mode limited to final results
- [ ] Show progress when verbose mode is enabled

## Tests

The parser and worker-pool tests pass with the race detector. Scanner-package
coverage is currently 88.8%.

- [x] Run tests with the race detector in GitHub Actions
- [ ] Test an open port with a temporary local TCP listener
- [ ] Test a closed port
- [x] Test port-list parsing
- [x] Test port-range parsing
- [x] Test invalid port input
- [x] Test duplicate-port handling
- [x] Test worker-pool result count
- [x] Test that worker-pool results are sorted
- [x] Test table output
- [x] Test JSON output
- [ ] Test CLI flag parsing
- [x] Add meaningful coverage for the main scanner functions

## Quality and automation

- [x] Format the current code with `gofmt`
- [x] Pass `go vet ./...`
- [x] Pass `golangci-lint`
- [x] Run `golangci-lint` in GitHub Actions
- [x] Run `go test -race -cover ./...` in GitHub Actions
- [x] Build a multi-stage Docker image
- [x] Run the container as a non-root user
- [x] Provide a Docker Compose service
- [x] Add real tests so the race job exercises concurrent code
- [ ] Consider adding `go vet ./...` as an explicit CI step

## Benchmark and documentation

- [ ] Implement a sequential scanning benchmark
- [ ] Implement a worker-pool scanning benchmark
- [ ] Compare both approaches using the same host and ports
- [ ] Record real results in `benchmark/results.md`
- [ ] Add real benchmark numbers to `README.md`
- [ ] Replace the draft README section with final installation instructions
- [ ] Document every CLI flag
- [ ] Add table-output and JSON-output examples
- [ ] Document local, Docker, and Docker Compose usage
- [ ] Document that users must only scan authorized targets

## Definition of done

- [ ] Scan 1,000 local ports with 100 workers in a few seconds
- [x] Pass meaningful tests with `go test -race ./...`
- [x] Pass `go vet ./...`
- [ ] Pass `golangci-lint run ./...`
- [x] Produce valid table and JSON output
- [ ] Complete the README with real examples and benchmark results
- [x] Maintain incremental Git history
- [x] Provide Docker, Docker Compose, and GitHub Actions configuration

## Recommended implementation order

1. Implement and test port parsing.
2. Implement and test the worker pool.
3. Replace the hard-coded `main.go` values with CLI flags.
4. Implement table and JSON output.
5. Add integration tests and increase coverage.
6. Run and document the benchmarks.
7. Finish the README and perform the definition-of-done checks.
