// This module is cross-platform (Windows and Linux).

package investigation_tools

import (
	"os"
	"path/filepath"
	"time"
)

// FileChange represents a file change event.
type FileChange struct {
	Path         string
	LastModified time.Time
}

// GetFileChangesByPath retrieves files modified within a specified time range.
//
// Description:
// - Walks through the directory and its subdirectories to find files modified within the time range.
// - Resolves symbolic links to their target paths.
//
// Parameters:
// - path (string): The directory path to search.
// - startDate (time.Time): The start of the time range.
// - endDate (time.Time): The end of the time range.
//
// Returns:
// - []FileChange: A slice of file change details.
// - error: An error if the directory cannot be accessed.
//
// Example Usage:
// ```go
// changes, err := GetFileChangesByPath("/path/to/directory", time.Now().Add(-24*time.Hour), time.Now())
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("File Changes:", changes)
//	}
//
// ```
func GetFileChangesByPath(path string, startDate, endDate time.Time) ([]FileChange, error) {
	var changes []FileChange

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle symbolic links
		if info.Mode()&os.ModeSymlink != 0 {
			resolvedPath, err := filepath.EvalSymlinks(filePath)
			if err != nil {
				return err
			}
			info, err = os.Stat(resolvedPath)
			if err != nil {
				return err
			}
			filePath = resolvedPath
		}

		// Check modification time
		if info.ModTime().After(startDate) && info.ModTime().Before(endDate) {
			changes = append(changes, FileChange{
				Path:         filePath,
				LastModified: info.ModTime(),
			})
		}
		return nil
	})

	return changes, err
}
