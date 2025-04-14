// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"os"
)

// GetAllDrives returns a list of all drives on the system.
//
// Description:
// - On Linux, it returns the root directory ("/").
// - On Windows, it iterates through all possible drive letters (A-Z) and checks if they exist.
//
// Parameters: None
//
// Returns:
// - []string: A slice of strings representing the available drives.
//
// Example Usage:
// ```go
// drives := GetAllDrives()
// fmt.Println("Available drives:", drives)
// ```
func GetAllDrives() []string {
	if os.PathSeparator == '/' {
		return []string{"/"} // Root directory for Linux
	}
	drives := []string{}
	for letter := 'A'; letter <= 'Z'; letter++ {
		drive := fmt.Sprintf("%c:\\", letter)
		if _, err := os.Stat(drive); err == nil {
			drives = append(drives, drive)
		}
	}
	return drives
}
