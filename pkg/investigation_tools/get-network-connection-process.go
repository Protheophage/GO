// This module is cross-platform (Windows and Linux).

package investigation_tools

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// NetworkConnection represents a network connection and its associated process.
type NetworkConnection struct {
	LocalAddress  string
	LocalPort     string
	RemoteAddress string
	RemotePort    string
	State         string
	ProcessName   string
	ProcessId     string
}

// GetNetworkConnectionProcess retrieves network connections and their associated processes.
//
// Description:
// - On Windows, uses `netstat` and `tasklist` to retrieve connection details.
// - On Linux, uses `ss` and `ps` to retrieve connection details.
//
// Parameters:
// - ipAddresses ([]string): A list of IP addresses to filter connections.
//
// Returns:
// - []NetworkConnection: A slice of network connections and their associated processes.
// - error: An error if the query fails.
//
// Example Usage:
// ```go
// connections, err := GetNetworkConnectionProcess([]string{"127.0.0.1"})
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Connections:", connections)
//	}
//
// ```
func GetNetworkConnectionProcess(ipAddresses []string) ([]NetworkConnection, error) {
	var results []NetworkConnection

	for _, ipAddress := range ipAddresses {
		log.Printf("Checking connections for IP address: %s", ipAddress)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("netstat", "-ano")
		} else {
			cmd = exec.Command("ss", "-tanp")
		}

		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute command: %v", err)
		}

		lines := strings.Split(string(output), "\n")
		found := false
		for _, line := range lines {
			fields := strings.Fields(line)
			if runtime.GOOS == "windows" {
				if len(fields) < 5 || !strings.Contains(fields[2], ipAddress) {
					continue
				}
			} else {
				if len(fields) < 6 || !strings.Contains(fields[4], ipAddress) {
					continue
				}
			}

			found = true
			var remoteAddr, state, processId string
			if runtime.GOOS == "windows" {
				remoteAddr = fields[2]
				state = fields[3]
				processId = fields[4]
			} else {
				remoteAddr = fields[4]
				state = fields[1]
				processId = strings.TrimPrefix(fields[5], "pid=")
			}

			remoteHost, remotePort, err := net.SplitHostPort(remoteAddr)
			if err != nil {
				log.Printf("Failed to parse remote address %s: %v", remoteAddr, err)
				continue
			}

			processName := "Unknown"
			if runtime.GOOS == "windows" {
				cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %s", processId))
				tasklistOutput, err := cmd.Output()
				if err == nil {
					tasklistLines := strings.Split(string(tasklistOutput), "\n")
					if len(tasklistLines) > 3 {
						processName = strings.Fields(tasklistLines[3])[0]
					}
				} else {
					log.Printf("Failed to retrieve process name for PID %s: %v", processId, err)
				}
			} else {
				cmd = exec.Command("ps", "-p", processId, "-o", "comm=")
				psOutput, err := cmd.Output()
				if err == nil {
					processName = strings.TrimSpace(string(psOutput))
				} else {
					log.Printf("Failed to retrieve process name for PID %s: %v", processId, err)
				}
			}

			results = append(results, NetworkConnection{
				LocalAddress:  fields[1],
				LocalPort:     strings.Split(fields[1], ":")[1],
				RemoteAddress: remoteHost,
				RemotePort:    remotePort,
				State:         state,
				ProcessName:   processName,
				ProcessId:     processId,
			})
		}

		if !found {
			log.Printf("No connections found for IP address: %s", ipAddress)
		}
	}

	if len(results) == 0 {
		log.Println("No network connections found for the specified IP addresses.")
	} else {
		log.Println("Network connections retrieved successfully.")
	}

	return results, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run get-network-connection-process.go <IP1> <IP2> ...")
		return
	}

	ipAddresses := os.Args[1:]
	results, err := GetNetworkConnectionProcess(ipAddresses)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print results
	for _, conn := range results {
		fmt.Printf("Local: %s:%s, Remote: %s:%s, State: %s, Process: %s (PID: %s)\n",
			conn.LocalAddress, conn.LocalPort, conn.RemoteAddress, conn.RemotePort, conn.State, conn.ProcessName, conn.ProcessId)
	}
}
