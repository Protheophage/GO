// This module is cross-platform (Windows and Linux).

package investigation_tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// WatchFileChangesByPath monitors a directory and its subdirectories for file changes.
//
// Description:
// - Uses the `fsnotify` package to watch for file creation, modification, deletion, and renaming events.
// - Normalizes paths for cross-platform compatibility.
//
// Parameters:
// - path (string): The directory path to monitor.
//
// Returns:
// - error: An error if the watcher fails to initialize or monitor the directory.
//
// Example Usage:
// ```go
// err := WatchFileChangesByPath("C:\\MyDirectory")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	}
//
// ```
func WatchFileChangesByPath(path string) error {
	// Normalize the path for cross-platform compatibility
	normalizedPath := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		normalizedPath = strings.ReplaceAll(normalizedPath, "/", "\\")
	} else {
		normalizedPath = strings.ReplaceAll(normalizedPath, "\\", "/")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	// Add the specified path and its subdirectories to the watcher
	err = filepath.Walk(normalizedPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path '%s': %v", p, err)
		}
		if info.IsDir() {
			if err := watcher.Add(p); err != nil {
				return fmt.Errorf("failed to add path '%s' to watcher: %v", p, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Monitoring changes in '%s' and its subdirectories. Press Ctrl+C to stop.", normalizedPath)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// Normalize event paths for cross-platform consistency
			eventName := filepath.Clean(event.Name)
			if runtime.GOOS == "windows" {
				eventName = strings.ReplaceAll(eventName, "/", "\\")
			} else {
				eventName = strings.ReplaceAll(eventName, "\\", "/")
			}

			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				log.Printf("File created: %s", eventName)
			case event.Op&fsnotify.Write == fsnotify.Write:
				log.Printf("File changed: %s", eventName)
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				log.Printf("File deleted: %s", eventName)
			case event.Op&fsnotify.Rename == fsnotify.Rename:
				log.Printf("File renamed: %s", eventName)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Error: %v", err)
		}
	}
}
