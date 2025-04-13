// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func SetProjectFolder(baseDir, taskDir string, changeDir bool) (string, error) {
	if baseDir == "" {
		if runtime.GOOS == "windows" {
			baseDir = filepath.Join(os.Getenv("SystemDrive"), "WorkDir")
		} else {
			baseDir = filepath.Join(os.Getenv("HOME"), "WorkDir")
		}
	}

	// Ensure base directory exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create base directory: %v", err)
		}
	}

	// Create task directory if specified
	finalPath := baseDir
	if taskDir != "" {
		finalPath = filepath.Join(baseDir, taskDir)
		if _, err := os.Stat(finalPath); os.IsNotExist(err) {
			if err := os.MkdirAll(finalPath, os.ModePerm); err != nil {
				return "", fmt.Errorf("failed to create task directory: %v", err)
			}
		}
	}

	// Change to the directory if requested
	if changeDir {
		if err := os.Chdir(finalPath); err != nil {
			return "", fmt.Errorf("failed to change directory: %v", err)
		}
	}

	return finalPath, nil
}
