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
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func labours(wg *sync.WaitGroup, task chan string, dialer net.Dialer, openPorts *[]string, mu *sync.Mutex, totalPorts *int, progress *int) {

	defer wg.Done()
	for addr := range task {

		conn, err := dialer.Dial("tcp", addr)
		if err == nil {

			mu.Lock()

			*openPorts = append(*openPorts, addr)
			*progress++

			mu.Unlock()

			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("No banner recieved from %s\n", addr)
			} else {
				fmt.Printf("Banner from %s: %s\n", addr, string(buf[:n]))
			}

			fmt.Printf("Connection to %s was successful\n", addr)
			conn.Close()
			mu.Lock() // Lock to safely modify the progress in a concurrent environment
			*progress++
			fmt.Printf("Scanning port %d/%d - Progress: %.2f%%\n", *progress, *totalPorts, float64(*progress)*100/float64(*totalPorts))
			mu.Unlock()
		} else {
			// Failed connection, print an error message
			fmt.Printf("Failed to connect to %s: %v\n", addr, err)
			mu.Lock() // Lock to safely modify the progress in a concurrent environment
			*progress++
			fmt.Printf("Scanning port %d/%d - Progress: %.2f%%\n", *progress, *totalPorts, float64(*progress)*100/float64(*totalPorts))
			mu.Unlock()
		}

	}
}

func scanTarget(target string, startPort int, endPort int, workers int, timeout int) {
	task := make(chan string, workers)
	var wg sync.WaitGroup

	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	var openPorts []string
	totalPorts := endPort - startPort + 1
	progress := 0
	mu := sync.Mutex{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go labours(&wg, task, dialer, &openPorts, &mu, &totalPorts, &progress)
	}

	// Distribute tasks for the current target
	for j := startPort; j <= endPort; j++ {
		port := strconv.Itoa(j)
		address := net.JoinHostPort(target, port)
		task <- address
	}
	close(task)

	// Wait for all workers to finish
	wg.Wait()

	duration := time.Since(time.Now())

	// Print results for the current target
	fmt.Printf("\nScan complete for target: %s\n", target)
	fmt.Printf("Open ports: %v\n", openPorts)
	fmt.Printf("Number of open ports: %d\n", len(openPorts))
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Total ports scanned: %d\n", totalPorts)

}

func main() {
	targets := flag.String("targets", "localhost", "Comma-separated list of target IPs or hostnames to scan")

	startPort := flag.Int("start-port", 1, "Starting port number")
	endPort := flag.Int("end-port", 1024, "Ending port number")

	workers := flag.Int("workers", 100, "Number of concurrent workers")

	timeout := flag.Int("timeout", 2, "Connection timeout in seconds")

	flag.Parse()

	// Split the targets into a slice
	targetList := strings.Split(*targets, ",")

	fmt.Printf("Scanning targets: %v\n", targetList)
	fmt.Printf("Scanning ports from %d to %d\n", *startPort, *endPort)
	fmt.Printf("Using %d concurrent workers\n", *workers)
	fmt.Printf("Connection timeout: %d seconds\n", *timeout)

	// Use goroutines to scan each target concurrently
	var wg sync.WaitGroup
	for _, target := range targetList {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			scanTarget(target, *startPort, *endPort, *workers, *timeout)
		}(target)
	}

	// Wait for all target scans to complete
	wg.Wait()

	fmt.Println("All scans complete!")
}
