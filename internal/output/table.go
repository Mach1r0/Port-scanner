package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/Mach1r0/port-scanner/internal/scanner"
)

func WriteTable(w io.Writer, results []scanner.Result) error {
	table := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(table, "PORT\tSTATUS\tDURATION"); err != nil {
		return fmt.Errorf("write table header: %w", err)
	}

	for _, result := range results {
		if _, err := fmt.Fprintf(
			table,
			"%d\t%s\t%s\n",
			result.Port,
			result.Status,
			result.Duration,
		); err != nil {
			return fmt.Errorf("write table result: %w", err)
		}
	}

	if err := table.Flush(); err != nil {
		return fmt.Errorf("flush table output: %w", err)
	}

	return nil
}
