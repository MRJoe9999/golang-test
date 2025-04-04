package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Information struct {
	Target     string   `json:"target"`
	OpenPorts  []string `json:"open_ports"`
	PortCount  int      `json:"port_count"`
	TimeTaken  string   `json:"time_taken"`
	TotalPorts int      `json:"total_ports"`
	Progress   float64  `json:"progress"`
}

func labours(wg *sync.WaitGroup, task chan string, dialer net.Dialer, openPorts *[]string, mu *sync.Mutex, totalPorts *int, progress *int, results *[]Information, target string) {
	defer wg.Done()

	for addr := range task {
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			*openPorts = append(*openPorts, addr)
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("No banner received from %s\n", addr)
			} else {
				fmt.Printf("Banner from %s: %s\n", addr, string(buf[:n]))
			}

			fmt.Printf("Connection to %s was successful\n", addr)
			conn.Close()

			mu.Lock()
			*progress++
			mu.Unlock()

			fmt.Printf("Scanning port %d/%d - Progress: %.2f%%\n", *progress, *totalPorts, float64(*progress)*100/float64(*totalPorts))
		} else {

			fmt.Printf("Failed to connect to %s: %v\n", addr, err)
			mu.Lock()
			*progress++
			mu.Unlock()

			fmt.Printf("Scanning port %d/%d - Progress: %.2f%%\n", *progress, *totalPorts, float64(*progress)*100/float64(*totalPorts))
		}
	}

	mu.Lock()
	result := Information{
		Target:     target,
		OpenPorts:  *openPorts,
		PortCount:  len(*openPorts),
		TotalPorts: *totalPorts,
		Progress:   float64(*progress) * 100 / float64(*totalPorts),
		TimeTaken:  time.Since(time.Now()).String(),
	}
	*results = append(*results, result)
	mu.Unlock()
}

func scanTarget(target string, ports []int, workers, timeout int, jsonOutput bool, results *[]Information) {
	task := make(chan string, workers)

	var wg sync.WaitGroup

	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	var openPorts []string
	totalPorts := len(ports)
	progress := 0
	mu := sync.Mutex{}

	// Starting the scan
	startTime := time.Now()

	// Create worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go labours(&wg, task, dialer, &openPorts, &mu, &totalPorts, &progress, results, target)
	}

	// Distribute the tasks (ports) to scan
	for _, port := range ports {
		portStr := strconv.Itoa(port)
		address := net.JoinHostPort(target, portStr)
		task <- address
	}
	close(task)

	// Wait for all workers to finish
	wg.Wait()

	// If the json flag is set, print the results in JSON format
	if jsonOutput {
		jsonResults, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		fmt.Println(string(jsonResults))
	} else {
		// Print results for the target in human-readable format
		duration := time.Since(startTime)
		fmt.Printf("\nScan complete for target: %s\n", target)
		fmt.Printf("Open ports: %v\n", openPorts)
		fmt.Printf("Number of open ports: %d\n", len(openPorts))
		fmt.Printf("Time taken: %v\n", duration)
		fmt.Printf("Total ports scanned: %d\n", totalPorts)
	}
}

func main() {
	targets := flag.String("targets", "localhost", "Comma-separated list of target IPs or hostnames to scan")

	portsStr := flag.String("ports", "", "Comma-separated list of specific ports to scan")

	workers := flag.Int("workers", 100, "Number of concurrent workers")

	timeout := flag.Int("timeout", 2, "Connection timeout in seconds")
	jsonOutput := flag.Bool("json", false, "Output scan results in JSON format")

	flag.Parse()

	targetList := strings.Split(*targets, ",")

	var ports []int
	if *portsStr != "" {
		for _, port := range strings.Split(*portsStr, ",") {
			p, err := strconv.Atoi(port)
			if err != nil {
				fmt.Printf("Invalid port number: %s\n", port)
				continue
			}
			ports = append(ports, p)
		}
	}

	if len(ports) == 0 {
		ports = append(ports, 1, 1024)
	}

	fmt.Printf("Scanning targets: %v\n", targetList)
	fmt.Printf("Scanning ports: %v\n", ports)
	fmt.Printf("Using %d concurrent workers\n", *workers)
	fmt.Printf("Connection timeout: %d seconds\n", *timeout)

	var results []Information

	var wg sync.WaitGroup
	for _, target := range targetList {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			scanTarget(target, ports, *workers, *timeout, *jsonOutput, &results)
		}(target)
	}

	wg.Wait()

	if !*jsonOutput {
		fmt.Println("All scans complete!")
	}

	fmt.Println("All scans complete!")
}
