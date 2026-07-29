package scanner

import (
	"net"
	"testing"
	"time"
)

const benchmarkPortCount = 1000

func BenchmarkSequentialScan(b *testing.B) {
	host, port := startBenchmarkServer(b)
	ports := repeatedPorts(port, benchmarkPortCount)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		for _, scanPort := range ports {
			ScanPort(host, scanPort, time.Second)
		}
	}
}

func BenchmarkWorkerPoolScan(b *testing.B) {
	host, port := startBenchmarkServer(b)
	ports := repeatedPorts(port, benchmarkPortCount)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := Scan(host, ports, 100, time.Second); err != nil {
			b.Fatalf("scan ports: %v", err)
		}
	}
}

func startBenchmarkServer(b *testing.B) (string, int) {
	b.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("start benchmark server: %v", err)
	}

	b.Cleanup(func() {
		_ = listener.Close()
	})

	go func() {
		for {
			connection, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}

			_ = connection.Close()
		}
	}()

	address := listener.Addr().(*net.TCPAddr)

	return address.IP.String(), address.Port
}

func repeatedPorts(port, count int) []int {
	ports := make([]int, count)
	for index := range ports {
		ports[index] = port
	}

	return ports
}
