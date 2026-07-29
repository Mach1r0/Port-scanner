package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Mach1r0/port-scanner/internal/scanner"
)

type jsonResult struct {
	Port     int            `json:"port"`
	Status   scanner.Status `json:"status"`
	Duration string         `json:"duration"`
}

func WriteJSON(w io.Writer, results []scanner.Result) error {
	jsonResults := make([]jsonResult, 0, len(results))

	for _, result := range results {
		jsonResults = append(jsonResults, jsonResult{
			Port:     result.Port,
			Status:   result.Status,
			Duration: result.Duration.String(),
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(jsonResults); err != nil {
		return fmt.Errorf("encode JSON results: %w", err)
	}

	return nil
}
