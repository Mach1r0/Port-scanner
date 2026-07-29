package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mach1r0/port-scanner/internal/output"
	"github.com/Mach1r0/port-scanner/internal/scanner"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	target := flag.String(
		"target",
		"",
		"IP address or hostname to scan",
	)

	portInput := flag.String(
		"ports",
		"1-1024",
		"ports to scan as a range (e.g., 1-1024) or a comma-separated list (e.g., 22,80,443)",
	)

	workerCount := flag.Int(
		"workers",
		100,
		"number of concurrent workers",
	)

	timeout := flag.Duration(
		"timeout",
		time.Second,
		"connection timeout",
	)

	outputFormat := flag.String(
		"output",
		"table",
		"output format: table or json",
	)

	verbose := flag.Bool(
		"verbose",
		false,
		"display scan progress",
	)

	flag.Parse()

	if strings.TrimSpace(*target) == "" {
		return errors.New("target is required")
	}

	ports, err := scanner.ParsePorts(*portInput)

	if err != nil {
		return fmt.Errorf("invalid ports: %w", err)
	}

	if *verbose {
		fmt.Fprintf(
			os.Stderr,
			"Scanning %s: %d ports with %d workers and timeout %s\n",
			*target,
			len(ports),
			*workerCount,
			*timeout,
		)
	}

	results, err := scanner.Scan(
		*target,
		ports,
		*workerCount,
		*timeout,
	)

	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	switch *outputFormat {
	case "table":
		if err := output.WriteTable(os.Stdout, results); err != nil {
			return fmt.Errorf("write table output: %w", err)
		}

	case "json":
		if err := output.WriteJSON(os.Stdout, results); err != nil {
			return fmt.Errorf("write JSON output: %w", err)
		}

	default:
		return fmt.Errorf(
			"unsupported output format %q: expected table or json",
			*outputFormat,
		)
	}

	return nil
}
