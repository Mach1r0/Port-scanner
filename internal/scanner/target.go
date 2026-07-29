package scanner

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func ValidateTarget(target string) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return errors.New("target is required")
	}

	if net.ParseIP(target) != nil {
		return nil
	}

	if strings.Contains(target, ":") || looksLikeIPv4(target) {
		return fmt.Errorf("invalid IP address %q", target)
	}

	hostname := strings.TrimSuffix(target, ".")
	if hostname == "" || len(hostname) > 253 {
		return fmt.Errorf("invalid hostname %q", target)
	}

	for _, label := range strings.Split(hostname, ".") {
		if !validHostnameLabel(label) {
			return fmt.Errorf("invalid hostname %q", target)
		}
	}

	return nil
}

func validHostnameLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 {
		return false
	}

	if label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}

	for _, character := range label {
		if character >= 'a' && character <= 'z' {
			continue
		}
		if character >= 'A' && character <= 'Z' {
			continue
		}
		if character >= '0' && character <= '9' {
			continue
		}
		if character == '-' {
			continue
		}

		return false
	}

	return true
}

func looksLikeIPv4(target string) bool {
	if strings.Count(target, ".") != 3 {
		return false
	}

	for _, character := range target {
		if character != '.' && (character < '0' || character > '9') {
			return false
		}
	}

	return true
}
