package main

import (
	"fmt"
	"time"

	"github.com/Mach1r0/port-scanner/internal/scanner"
)

func main() {
	result := scanner.ScanPort("localhost", 80, 2*time.Second)
	fmt.Printf("Port: %d, Status: %s, Duration: %s\n",
		result.Port, result.Status, result.Duration)
}
