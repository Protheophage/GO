// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// InstallAppFromWeb downloads and installs an application from a given URL.
//
// Description:
// - Checks if the application is already installed by verifying the existence of a specified path.
// - Downloads the installer from the provided URL and executes it with optional arguments.
// - Cleans up the installer file after installation.
//
// Parameters:
// - installCheckPath (string): Path to check if the application is already installed.
// - installerURL (string): URL to download the installer.
// - appName (string): Name of the application (used for naming the installer file).
// - installerArgs (string): Arguments to pass to the installer (default is "/S" for silent installation).
//
// Returns:
// - error: An error if the installation fails.
//
// Example Usage:
// ```go
// err := InstallAppFromWeb("C:\\Program Files\\ExampleApp", "https://example.com/installer.exe", "ExampleApp", "/S")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Application installed successfully.")
//	}
//
// ```
func InstallAppFromWeb(installCheckPath, installerURL, appName, installerArgs string) error {
	if installerArgs == "" {
		installerArgs = "/S"
	}

	// Check if the application is already installed
	if installCheckPath != "" {
		if _, err := os.Stat(installCheckPath); err == nil {
			fmt.Println("Application is already installed. Skipping installation.")
			return nil
		}
	}

	// Download the installer
	installerPath := filepath.Join(os.TempDir(), appName)
	resp, err := http.Get(installerURL)
	if err != nil {
		return fmt.Errorf("failed to download the installer: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(installerPath)
	if err != nil {
		return fmt.Errorf("failed to create installer file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save the installer: %v", err)
	}

	// Run the installer silently
	cmd := exec.Command(installerPath, installerArgs)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run the installer: %v", err)
	}

	// Clean up the installer file
	if err := os.Remove(installerPath); err != nil {
		return fmt.Errorf("failed to remove the installer file: %v", err)
	}

	return nil
}
