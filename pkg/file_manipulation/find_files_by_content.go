package file_manipulation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// FindFilesByContent searches for files containing specific content.
func FindFilesByContent(stringToFind, fileTypeToSearch string, maxFileSizeKB int, searchAllDrives bool, checkThisDisk string) ([]string, error) {
	if checkThisDisk == "" {
		if os.PathSeparator == '/' {
			checkThisDisk = "/"
		} else {
			checkThisDisk = filepath.Join(os.Getenv("SystemDrive"), "")
		}
	}

	var drives []string
	if searchAllDrives {
		drives = random_utilities.GetAllDrives()
	} else {
		drives = []string{checkThisDisk}
	}

	var foundFiles []string
	for _, drive := range drives {
		fmt.Printf("Searching for content in: %s\n", drive)
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if filepath.Ext(path) == fileTypeToSearch && info.Size() <= int64(maxFileSizeKB*1024) {
				file, err := os.Open(path)
				if err != nil {
					return nil
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), stringToFind) {
						foundFiles = append(foundFiles, path)
						break
					}
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error processing drive %s: %v", drive, err)
		}
	}

	return foundFiles, nil
}
