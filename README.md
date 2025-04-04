This is Part-B of the test, this is also a tool to check aivailable ports for a hostname.
This is have extra features:
- Banner reading - this tries to read the response from a server
- Progress indicator - this feature shows progress of each all the ports being scanned
- Scan multiple targets - accepts more than one target to scan at the same time 
- readability into json format - formats the results into a json format for better readability
- Specific port - allow user to enter specific ports to scan. 

#  go run main.go -targets=scanme.nmap.org -ports=22 -json
# go run main.go -targets=scanme.nmap.org -ports=22,80,443 -json
# go run main.go -targets=example.com,scanme.nmap.org -ports=22,80,443 -workers=10

# example of output

``` Scanning targets: [example.com scanme.nmap.org]
Scanning ports: [22 80 443]
Using 10 concurrent workers
Connection timeout: 2 seconds
Banner from scanme.nmap.org:22: SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13

Connection to scanme.nmap.org:22 was successful
Scanning port 1/3 - Progress: 33.33%
Failed to connect to example.com:22: dial tcp 23.215.0.138:22: i/o timeout
Scanning port 1/3 - Progress: 33.33%
Failed to connect to scanme.nmap.org:443: dial tcp 45.33.32.156:443: i/o timeout
Scanning port 2/3 - Progress: 66.67%
No banner received from scanme.nmap.org:80
Connection to scanme.nmap.org:80 was successful
Scanning port 3/3 - Progress: 100
Scan complete for target: scanme.nmap.org
Open ports: [scanme.nmap.org:80 scanme.nmap.org:22]
Number of open ports: 2
Time taken: 2.246895126s
Total ports scanned: 3
No banner received from example.com:443
Connection to example.com:443 was successful
Scanning port 2/3 - Progress: 66.67%
No banner received from example.com:80
Connection to example.com:80 was successful
Scanning port 3/3 - Progress: 100.00%
Scan complete for target: example.com
Open ports: [example.com:443 example.com:80]
Number of open ports: 2
Time taken: 2.258699674s
Total ports scanned: 3
All scans complete!
All scans complete!

# clone repository
git clone -b <part-B> <https://github.com/MRJoe9999/golang-test.git>