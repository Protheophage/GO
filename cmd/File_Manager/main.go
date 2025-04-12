package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Protheophage/GO/pkg/file_manipulation"
)

func main() {
	// Define subcommands
	countCmd := flag.NewFlagSet("count", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	contentCmd := flag.NewFlagSet("content", flag.ExitOnError)
	extensionCmd := flag.NewFlagSet("extension", flag.ExitOnError)

	// Common flags
	searchAllDrives := flag.Bool("all", false, "Search all drives")
	checkThisDisk := flag.String("disk", "", "Specific disk to search")

	// Flags for count
	countPattern := countCmd.String("pattern", "*", "File pattern to count")

	// Flags for remove
	removePattern := removeCmd.String("pattern", "*", "File pattern to remove")

	// Flags for find
	findPattern := findCmd.String("pattern", "*", "File pattern to find")

	// Flags for content
	contentString := contentCmd.String("string", "", "String to search for in files")
	contentType := contentCmd.String("type", "", "File type to search")
	contentMaxSize := contentCmd.Int("maxsize", 1024, "Max file size in KB")

	// Flags for extension
	extensionPattern := extensionCmd.String("pattern", "*", "File pattern to change extension")
	newExtension := extensionCmd.String("new", ".txt", "New file extension")

	// Parse subcommands
	if len(os.Args) < 2 {
		fmt.Println("File Manager CLI Application")
		fmt.Println("Usage:")
		fmt.Println("  file_manager <command> [flags]")
		fmt.Println("Commands:")
		fmt.Println("  count      Count files matching a pattern")
		fmt.Println("  remove     Remove files matching a pattern")
		fmt.Println("  find       Find files matching a pattern")
		fmt.Println("  content    Find files containing specific content")
		fmt.Println("  extension  Change file extensions")
		fmt.Println("Use 'file_manager <command> -h' for more information about a command.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "count":
		countCmd.Usage = func() {
			fmt.Println("Usage: file_manager count [flags]")
			fmt.Println("Flags:")
			countCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file_manager count -pattern=\"*.txt\" -all")
		}
		countCmd.Parse(os.Args[2:])
		if *countPattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		count, err := file_manipulation.GetFilesCount(*countPattern, *searchAllDrives, *checkThisDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Found %d files matching pattern '%s'\n", count, *countPattern)

	case "remove":
		removeCmd.Usage = func() {
			fmt.Println("Usage: file_manager remove [flags]")
			fmt.Println("Flags:")
			removeCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file_manager remove -pattern=\"*.log\" -disk=\"C:\\\"")
		}
		removeCmd.Parse(os.Args[2:])
		if *removePattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		if err := file_manipulation.RemoveFiles(*removePattern, *searchAllDrives, *checkThisDisk); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Files removed successfully.")

	case "find":
		findCmd.Usage = func() {
			fmt.Println("Usage: file_manager find [flags]")
			fmt.Println("Flags:")
			findCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file_manager find -pattern=\"*.go\" -disk=\"/\"")
		}
		findCmd.Parse(os.Args[2:])
		if *findPattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		files, err := file_manipulation.FindFiles(*findPattern, *searchAllDrives, *checkThisDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Found files:", files)

	case "content":
		contentCmd.Usage = func() {
			fmt.Println("Usage: file_manager content [flags]")
			fmt.Println("Flags:")
			contentCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file_manager content -string=\"TODO\" -type=\".go\" -maxsize=512 -all")
		}
		contentCmd.Parse(os.Args[2:])
		if *contentString == "" {
			fmt.Println("Error: Search string cannot be empty.")
			os.Exit(1)
		}
		if *contentType == "" {
			fmt.Println("Error: File type cannot be empty.")
			os.Exit(1)
		}
		if *contentMaxSize <= 0 {
			fmt.Println("Error: Max file size must be greater than 0.")
			os.Exit(1)
		}
		files, err := file_manipulation.FindFilesByContent(*contentString, *contentType, *contentMaxSize, *searchAllDrives, *checkThisDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Found files containing string:", files)

	case "extension":
		extensionCmd.Usage = func() {
			fmt.Println("Usage: file_manager extension [flags]")
			fmt.Println("Flags:")
			extensionCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file_manager extension -pattern=\"*.txt\" -new=\".md\" -disk=\"/\"")
		}
		extensionCmd.Parse(os.Args[2:])
		if *extensionPattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		if *newExtension == "" {
			fmt.Println("Error: New extension cannot be empty.")
			os.Exit(1)
		}
		if err := file_manipulation.SetFilesExtension(*extensionPattern, *newExtension, *searchAllDrives, *checkThisDisk); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("File extensions updated successfully.")

	default:
		fmt.Println("Unknown command. Use 'file_manager -h' for help.")
		os.Exit(1)
	}
}
