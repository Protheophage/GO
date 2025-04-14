// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// SetProjectFolder creates a project folder structure and optionally changes the working directory.
//
// Description:
// - Creates a base directory and an optional task directory inside it.
// - Changes the working directory to the final path if requested.
//
// Parameters:
// - baseDir (string): The base directory path. Defaults to "WorkDir" in the system drive or home directory.
// - taskDir (string): The name of the task directory to create inside the base directory.
// - changeDir (bool): Whether to change the working directory to the final path.
//
// Returns:
// - string: The final path of the created directory.
// - error: An error if directory creation or changing the working directory fails.
//
// Example Usage:
// ```go
// path, err := SetProjectFolder("", "Task1", true)
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Project folder set to:", path)
//	}
//
// ```
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
