package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Mach1r0/port-scanner/internal/output"
	"github.com/Mach1r0/port-scanner/internal/scanner"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("scanner", flag.ContinueOnError)
	flags.SetOutput(stderr)

	target := flags.String(
		"target",
		"",
		"IP address or hostname to scan",
	)
	portInput := flags.String(
		"ports",
		"1-1024",
		"ports to scan as a range (e.g., 1-1024) or comma-separated values (e.g., 22,80,443)",
	)
	workerCount := flags.Int(
		"workers",
		100,
		"number of concurrent workers",
	)
	timeout := flags.Duration(
		"timeout",
		time.Second,
		"connection timeout",
	)
	outputFormat := flags.String(
		"output",
		"table",
		"output format: table or json",
	)
	verbose := flags.Bool(
		"verbose",
		false,
		"display scan progress",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	if err := scanner.ValidateTarget(*target); err != nil {
		return err
	}
	if *workerCount <= 0 {
		return errors.New("workers must be greater than zero")
	}
	if *timeout <= 0 {
		return errors.New("timeout must be greater than zero")
	}
	if *outputFormat != "table" && *outputFormat != "json" {
		return fmt.Errorf(
			"unsupported output format %q: expected table or json",
			*outputFormat,
		)
	}

	ports, err := scanner.ParsePorts(*portInput)
	if err != nil {
		return fmt.Errorf("invalid ports: %w", err)
	}

	var progress scanner.ProgressFunc
	var progressErr error
	if *verbose {
		if _, err := fmt.Fprintf(
			stderr,
			"Scanning %s: %d ports with %d workers and timeout %s\n",
			*target,
			len(ports),
			*workerCount,
			*timeout,
		); err != nil {
			return fmt.Errorf("write scan summary: %w", err)
		}

		progress = func(completed, total int, result scanner.Result) {
			if progressErr != nil {
				return
			}

			_, progressErr = fmt.Fprintf(
				stderr,
				"[%d/%d] port %d: %s\n",
				completed,
				total,
				result.Port,
				result.Status,
			)
		}
	}

	results, err := scanner.ScanWithProgress(
		*target,
		ports,
		*workerCount,
		*timeout,
		progress,
	)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	if progressErr != nil {
		return fmt.Errorf("write scan progress: %w", progressErr)
	}

	switch *outputFormat {
	case "table":
		if err := output.WriteTable(stdout, results); err != nil {
			return fmt.Errorf("write table output: %w", err)
		}
	case "json":
		if err := output.WriteJSON(stdout, results); err != nil {
			return fmt.Errorf("write JSON output: %w", err)
		}
	}

	return nil
}
