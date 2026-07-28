package scanner

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

func ParsePorts(portStr string) ([]int, error) {
	input := strings.TrimSpace(portStr)
	if input == "" {
		return nil, fmt.Errorf("ports cannot be empty")
	}

	seen := make(map[int]struct{})

	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("invalid empty port")
		}

		if strings.Contains(part, "-") {
			if err := addPortRange(part, seen); err != nil {
				return nil, err
			}
			continue
		}

		port, err := parsePort(part)
		if err != nil {
			return nil, err
		}

		seen[port] = struct{}{}
	}

	ports := make([]int, 0, len(seen))
	for port := range seen {
		ports = append(ports, port)
	}

	sort.Ints(ports)

	return ports, nil
}

func addPortRange(input string, seen map[int]struct{}) error {
	startText, endText, found := strings.Cut(input, "-")
	if !found || strings.Contains(endText, "-") {
		return fmt.Errorf("invalid port range %q", input)
	}

	start, err := parsePort(strings.TrimSpace(startText))
	if err != nil {
		return err
	}

	end, err := parsePort(strings.TrimSpace(endText))
	if err != nil {
		return err
	}

	if start > end {
		return fmt.Errorf("range start %d is greater than end %d", start, end)
	}

	for port := start; port <= end; port++ {
		seen[port] = struct{}{}
	}

	return nil
}

func parsePort(input string) (int, error) {
	port, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q", input)
	}

	if port < minPort || port > maxPort {
		return 0, fmt.Errorf(
			"port %d is outside the valid range %d-%d",
			port,
			minPort,
			maxPort,
		)
	}

	return port, nil
}
