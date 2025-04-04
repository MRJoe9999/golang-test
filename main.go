// Filename: main.go
// NAME: Jose Teck
package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func labours(wg *sync.WaitGroup, task chan string, dialer net.Dialer) {

	defer wg.Done()
	for addr := range task {

		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			fmt.Printf("Connection to %s was successful\n", addr)
			conn.Close()
		} else {
			// Failed connection, print an error message
			fmt.Printf("Failed to connect to %s: %v\n", addr, err)
		}
	}
}

func main() {
	target := flag.String("target", "localhost", "Target IP or hostname")

	startPort := flag.Int("start-port", 1, "Starting port number")
	endPort := flag.Int("end-port", 1024, "Ending port number")

	workers := flag.Int("workers", 100, "Number of concurrent workers")

	timeout := flag.Int("timeout", 2, "Connectione timout in seconds")

	flag.Parse()

	fmt.Println("Scanning target:", *target)
	fmt.Printf("Scanning ports from %d to %d\n", *startPort, *endPort)
	fmt.Printf("Using %d concurrent workers\n", *workers)
	fmt.Printf("Connection timeout: %d seconds\n", *timeout)

	task := make(chan string, *workers)

	var wg sync.WaitGroup

	dialer := net.Dialer{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go labours(&wg, task, dialer)
	}

	for j := *startPort; j <= *endPort; j++ {
		port := strconv.Itoa(j)
		address := net.JoinHostPort(*target, port)
		task <- address
	}

	close(task)
	wg.Wait()

	fmt.Println("Scan complete!")

}
