package scanner

import (
	"errors"
	"sort"
	"sync"
	"time"
)

func Scan(
	host string,
	ports []int,
	workerCount int,
	timeout time.Duration,
) ([]Result, error) {
	if workerCount <= 0 {
		return nil, errors.New("worker count must be greater than zero")
	}

	if len(ports) == 0 {
		return []Result{}, nil
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
