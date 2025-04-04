// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

/*
package main

import (

	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

)

	func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer) {
		defer wg.Done()
		maxRetries := 3
		for addr := range tasks {
			var success bool
			for i := range maxRetries {
				conn, err := dialer.Dial("tcp", addr)
				if err == nil {
					conn.Close()
					fmt.Printf("Connection to %s was successful\n", addr)
					success = true
					break
				}
				backoff := time.Duration(1<<i) * time.Second
				fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1, addr, backoff)
				time.Sleep(backoff)
			}
			if !success {
				fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
			}
		}
	}

func main() {

		var wg sync.WaitGroup
		tasks := make(chan string, 100)

		target := "scanme.nmap.org"

		dialer := net.Dialer{
			Timeout: 5 * time.Second,
		}

		workers := 100

		for i := 1; i <= workers; i++ {
			wg.Add(1)
			go worker(&wg, tasks, dialer)
		}

		ports := 512

		for p := 1; p <= ports; p++ {
			port := strconv.Itoa(p)
			address := net.JoinHostPort(target, port)
			tasks <- address
		}
		close(tasks)
		wg.Wait()
	}
*/
/*package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPort(target string, port int, timeout time.Duration, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil {
		conn.Close()
		results <- port
	}
}

func main() {
	var target string
	var startPort, endPort, workers int
	var timeoutSec int

	flag.StringVar(&target, "target", "localhost", "Target IP or hostname")
	flag.IntVar(&startPort, "start-port", 1, "Start of port range")
	flag.IntVar(&endPort, "end-port", 1024, "End of port range")
	flag.IntVar(&workers, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&timeoutSec, "timeout", 2, "Connection timeout in seconds")
	flag.Parse()

	timeout := time.Duration(timeoutSec) * time.Second
	ports := make(chan int, workers)
	results := make(chan int)
	var openPorts []int
	var wg sync.WaitGroup
	startTime := time.Now()

	// Worker pool
	for i := 0; i < workers; i++ {
		go func() {
			for port := range ports {
				scanPort(target, port, timeout, results, &wg)
			}
		}()
	}

	// Sending ports to be scanned
	go func() {
		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			ports <- port
		}
		close(ports)
	}()

	// Collect results
	go func() {
		for port := range results {
			openPorts = append(openPorts, port)
		}
	}()

	wg.Wait()
	close(results)
	elapsedTime := time.Since(startTime)

	fmt.Println("Scan complete!")
	fmt.Printf("Open ports: %v\n", openPorts)
	fmt.Printf("Time taken: %s\n", elapsedTime)
	fmt.Printf("Total ports scanned: %d\n", endPort-startPort+1)
}
*/

package main

import (
	"flag"
	"fmt"
)

func main() {
	target := flag.String("target", "localhost", "Target IP or hostname")

	startPort := flag.Int("start-port", 1, "Starting port number")
	endPort := flag.Int("end-port", 1024, "Ending port number")

	workers := flag.Int("workers", 100, "Number of concurrent workers")

	flag.Parse()

	fmt.Println("Scanning target:", *target)
	fmt.Printf("Scanning ports from %d to %d\n", *startPort, *endPort)
	fmt.Printf("Using %d concurrent workers\n", *workers)
}
