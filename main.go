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

// function for workers to connect to the target address
// and port
// it takes a wait group, a channel for tasks, and a dialer
// for establishing connections
// it will attempt to connect to the address and port
// if successful, it will close the connection
// if failed, it will print an error message
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

// main function to parse command line arguments
// and set up the scanning process
// it takes the target address, starting and ending port numbers,
// number of workers, and connection timeout
// it creates a channel for tasks and a wait group for synchronization
// it creates a dialer with the specified timeout
// it starts the workers and sends the tasks to the channel
// it waits for all workers to finish and closes the channel
// it prints a message when the scan is complete
// it uses the flag package to parse command line arguments
// it uses the net package to create a dialer and establish connections

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
