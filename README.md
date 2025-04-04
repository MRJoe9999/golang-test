This is part A of the Test. This tool is use to target specific targets example the user can enter which hostname they want to check their ports availability. This one here lets the user enter the start port and end port. Also same goes to the workers, the one that checks the ports, the user can enter how much workers  and ofcourse the timeout option which will dictate how much time each port will take to check their availability. 
And finnally it prints a summary of the what all ports are available and how much ports are completely scanned. 
Example of cammand : 
# go run . -target scanme.nmap.org -start-port 1 -end-port 1000 -workers 500 -timeout 5
specify the target, the start port and end port, workers, and timeout. 
Here are the first five lines of output expected 
# Scanning target: scanme.nmap.org
# Scanning ports from 1 to 1000
# Using 500 concurrent workers
# Connection timeout: 5 seconds
# Connection to scanme.nmap.org:22 was successful

and then the rest are mostly failed like this 
# Failed to connect to scanme.nmap.org:981: dial tcp 45.33.32.156:981: i/o timeout
and then 
# Scan complete!

clone this repository 
git clone -b <main> <https://github.com/MRJoe9999/golang-test.git>
