package scanner

import (
	"net"
	"strconv"
	"time"
)

type Status string

const (
	StatusOpen     Status = "open"
	StatusClosed   Status = "closed"
	StatusFiltered Status = "filtered"
)

type Result struct {
	Port     int
	Status   Status
	Duration time.Duration
}

func ScanPort(host string, port int, timeout time.Duration) Result {
	startedAt := time.Now()
	address := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(startedAt)

	if err == nil {
		defer func() {
			_ = conn.Close()
		}()

		return Result{
			Port:     port,
			Status:   StatusOpen,
			Duration: duration,
		}
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return Result{
			Port:     port,
			Status:   StatusFiltered,
			Duration: duration,
		}
	}

	return Result{
		Port:     port,
		Status:   StatusClosed,
		Duration: duration,
	}
}
