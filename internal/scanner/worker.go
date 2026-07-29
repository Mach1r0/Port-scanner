package scanner

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type ProgressFunc func(completed, total int, result Result)

func Scan(
	host string,
	ports []int,
	workerCount int,
	timeout time.Duration,
) ([]Result, error) {
	return ScanWithProgress(host, ports, workerCount, timeout, nil)
}

func ScanWithProgress(
	host string,
	ports []int,
	workerCount int,
	timeout time.Duration,
	progress ProgressFunc,
) ([]Result, error) {
	if err := ValidateTarget(host); err != nil {
		return nil, err
	}

	if workerCount <= 0 {
		return nil, errors.New("worker count must be greater than zero")
	}

	if timeout <= 0 {
		return nil, errors.New("timeout must be greater than zero")
	}

	if len(ports) == 0 {
		return []Result{}, nil
	}

	for _, port := range ports {
		if port < minPort || port > maxPort {
			return nil, fmt.Errorf(
				"port %d is outside the valid range %d-%d",
				port,
				minPort,
				maxPort,
			)
		}
	}

	if workerCount > len(ports) {
		workerCount = len(ports)
	}

	jobs := make(chan int)
	results := make(chan Result)

	var workers sync.WaitGroup
	workers.Add(workerCount)

	for range workerCount {
		go scanWorker(host, timeout, jobs, results, &workers)
	}

	go func() {
		defer close(jobs)

		for _, port := range ports {
			jobs <- port
		}
	}()

	go func() {
		workers.Wait()
		close(results)
	}()

	scanResult := make([]Result, 0, len(ports))

	for result := range results {
		scanResult = append(scanResult, result)
		if progress != nil {
			progress(len(scanResult), len(ports), result)
		}
	}

	sort.Slice(scanResult, func(i, j int) bool {
		return scanResult[i].Port < scanResult[j].Port
	})

	return scanResult, nil
}

func scanWorker(
	host string,
	timeout time.Duration,
	jobs <-chan int,
	results chan<- Result,
	workers *sync.WaitGroup,
) {
	defer workers.Done()

	for port := range jobs {
		results <- ScanPort(host, port, timeout)
	}
}
