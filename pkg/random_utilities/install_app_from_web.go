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
