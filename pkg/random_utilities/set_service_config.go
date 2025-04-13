// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"os/exec"
	"runtime"
)

func SetServiceConfig(serviceName, recover, status, startup string) error {
	if recover == "" {
		recover = "restart"
	}
	if status == "" {
		status = "start"
	}
	if startup == "" {
		startup = "automatic"
	}

	var action string
	switch recover {
	case "restart":
		action = "restart/30000/restart/30000/restart/30000"
	case "noaction":
		action = "//////"
	case "reboot":
		action = "reboot/30000"
	default:
		return fmt.Errorf("invalid recover option: %s", recover)
	}

	if runtime.GOOS == "windows" {
		// Windows-specific commands
		if err := exec.Command("sc.exe", "config", serviceName, "start=", startup).Run(); err != nil {
			return fmt.Errorf("failed to set startup type: %v", err)
		}
		if err := exec.Command("sc.exe", "failure", serviceName, "actions=", action, "reset=", "4000").Run(); err != nil {
			return fmt.Errorf("failed to set recovery actions: %v", err)
		}
		if status == "start" {
			if err := exec.Command("net", "start", serviceName).Run(); err != nil {
				return fmt.Errorf("failed to start service: %v", err)
			}
		} else if status == "stop" {
			if err := exec.Command("net", "stop", serviceName).Run(); err != nil {
				return fmt.Errorf("failed to stop service: %v", err)
			}
		}
	} else if runtime.GOOS == "linux" {
		// Linux-specific commands
		if err := exec.Command("systemctl", "enable", serviceName).Run(); err != nil && startup == "automatic" {
			return fmt.Errorf("failed to enable service: %v", err)
		}
		if err := exec.Command("systemctl", "disable", serviceName).Run(); err != nil && startup == "manual" {
			return fmt.Errorf("failed to disable service: %v", err)
		}
		if recover == "restart" {
			if err := exec.Command("systemctl", "set-property", serviceName, "Restart=on-failure", "RestartSec=30").Run(); err != nil {
				return fmt.Errorf("failed to set recovery actions: %v", err)
			}
		}
		if status == "start" {
			if err := exec.Command("systemctl", "start", serviceName).Run(); err != nil {
				return fmt.Errorf("failed to start service: %v", err)
			}
		} else if status == "stop" {
			if err := exec.Command("systemctl", "stop", serviceName).Run(); err != nil {
				return fmt.Errorf("failed to stop service: %v", err)
			}
		}
	} else {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	fmt.Printf("Service %s configured with startup=%s, recover=%s, status=%s\n", serviceName, startup, recover, status)
	return nil
}
