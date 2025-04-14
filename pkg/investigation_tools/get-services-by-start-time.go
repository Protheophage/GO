// This module is cross-platform (Windows and Linux).

package investigation_tools

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Service represents a system service.
type Service struct {
	Name        string
	DisplayName string
	StartTime   time.Time
}

// GetServicesByStartTime retrieves services that started within a specified time range.
//
// Description:
// - On Windows, uses PowerShell to query service details.
// - On Linux, uses `systemctl` and `ps` to retrieve service information.
//
// Parameters:
// - serviceNames ([]string): A list of service names to filter (use "*" for all services).
// - startDate (time.Time): The start of the time range.
// - endDate (time.Time): The end of the time range.
//
// Returns:
// - []Service: A slice of services matching the criteria.
// - error: An error if the query fails.
//
// Example Usage:
// ```go
// services, err := GetServicesByStartTime([]string{"MyService"}, time.Now().Add(-24*time.Hour), time.Now())
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Services:", services)
//	}
//
// ```
func GetServicesByStartTime(serviceNames []string, startDate, endDate time.Time) ([]Service, error) {
	// Set default time range if not provided
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour) // Beginning of yesterday
	}
	if endDate.IsZero() {
		endDate = time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Second) // End of today
	}

	if len(serviceNames) == 0 {
		serviceNames = []string{"*"}
	}

	var services []Service
	var err error

	switch runtime.GOOS {
	case "windows":
		services, err = getWindowsServices(serviceNames, startDate, endDate)
	case "linux":
		services, err = getLinuxServices(serviceNames, startDate, endDate)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err != nil {
		return nil, err
	}

	return services, nil
}

func getWindowsServices(serviceNames []string, startDate, endDate time.Time) ([]Service, error) {
	cmd := exec.Command("powershell", "-Command", "Get-CimInstance -ClassName Win32_Service | Select-Object Name, DisplayName, ProcessId")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var services []Service

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		name := fields[0]
		displayName := fields[1]
		processID := fields[2]

		// Check if the service name matches the filter
		matches := false
		for _, serviceName := range serviceNames {
			if serviceName == "*" || strings.Contains(name, serviceName) {
				matches = true
				break
			}
		}
		if !matches {
			continue
		}

		// Get process start time
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("(Get-Process -Id %s).StartTime.ToString('yyyy-MM-dd HH:mm:ss')", processID))
		startTimeOutput, err := cmd.Output()
		if err != nil {
			continue
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(string(startTimeOutput)))
		if err != nil || startTime.Before(startDate) || startTime.After(endDate) {
			continue
		}

		services = append(services, Service{
			Name:        name,
			DisplayName: displayName,
			StartTime:   startTime,
		})
	}

	return services, nil
}

func getLinuxServices(serviceNames []string, startDate, endDate time.Time) ([]Service, error) {
	cmd := exec.Command("bash", "-c", "systemctl list-units --type=service --no-pager --all")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var services []Service

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		name := fields[0]
		displayName := fields[0] // Linux services typically don't have a separate display name
		status := fields[3]

		// Check if the service name matches the filter
		matches := false
		for _, serviceName := range serviceNames {
			if serviceName == "*" || strings.Contains(name, serviceName) {
				matches = true
				break
			}
		}
		if !matches || status != "running" {
			continue
		}

		// Get process start time
		cmd := exec.Command("bash", "-c", fmt.Sprintf("ps -o lstart= -p $(systemctl show -p MainPID --value %s)", name))
		startTimeOutput, err := cmd.Output()
		if err != nil {
			continue
		}

		startTime, err := time.Parse("Mon Jan 2 15:04:05 2006", strings.TrimSpace(string(startTimeOutput)))
		if err != nil || startTime.Before(startDate) || startTime.After(endDate) {
			continue
		}

		services = append(services, Service{
			Name:        name,
			DisplayName: displayName,
			StartTime:   startTime,
		})
	}

	return services, nil
}
