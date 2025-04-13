package main

import (
	"fmt"
	"github.com/Protheophage/GO/pkg/random_utilities"
)

func main() {
	// Test case 1: Valid download
	err := random_utilities.GetFileFromWeb(
		"https://winscp.net/download/WinSCP-6.5-Setup.exe/download", // Replace with a valid URL
		"D:\\Test\\WinSCP-6.5-Setup.exe",
		true,
	)
	if err != nil {
		fmt.Printf("Test case 1 failed: %v\n", err)
	} else {
		fmt.Println("Test case 1 passed: File downloaded successfully.")
	}

	// Test case 2: File already exists, overwrite disabled
	err = random_utilities.GetFileFromWeb(
		"https://winscp.net/download/WinSCP-6.5-Setup.exe/download", // Replace with a valid URL
		"D:\\Test\\WinSCP-6.5-Setup.exe",
		false,
	)
	if err != nil {
		fmt.Printf("Test case 2 passed: %v\n", err)
	} else {
		fmt.Println("Test case 2 failed: File should not have been overwritten.")
	}

	// Test case 3: Invalid URL
	err = random_utilities.GetFileFromWeb(
		"invalid-url",
		"D:\\Test\\WinSCP-6.5-Setup.exe",
		true,
	)
	if err != nil {
		fmt.Printf("Test case 3 passed: %v\n", err)
	} else {
		fmt.Println("Test case 3 failed: Invalid URL should have caused an error.")
	}

	// Test case 4: Non-existent directory
	err = random_utilities.GetFileFromWeb(
		"https://winscp.net/download/WinSCP-6.5-Setup.exe/download", // Replace with a valid URL
		"D:\\Test\\WinSCP-6.5-Setup.exe",
		true,
	)
	if err != nil {
		fmt.Printf("Test case 4 failed: %v\n", err)
	} else {
		fmt.Println("Test case 4 passed: File downloaded to a new directory.")
	}
}
