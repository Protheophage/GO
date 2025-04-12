package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Protheophage/GO/pkg/file_manipulation"
)

func main() {
	// Global help flag
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("File Manager CLI Application")
		fmt.Println("Usage:")
		fmt.Println("  file-manager <command> [flags]")
		fmt.Println("Commands:")
		fmt.Println("  count      Count files matching a pattern")
		fmt.Println("    Flags:")
		fmt.Println("      -pattern: File pattern to count (e.g., '*.txt').")
		fmt.Println("      -all: Search all drives (default: false).")
		fmt.Println("      -disk: Specify a disk to search.")
		fmt.Println("             Windows: 'C:\\' or 'D:\\'")
		fmt.Println("             Linux: '/' or '/home/user/'")
		fmt.Println()
		fmt.Println("  remove     Remove files matching a pattern")
		fmt.Println("    Flags:")
		fmt.Println("      -pattern: File pattern to remove (e.g., '*.log').")
		fmt.Println("      -all: Search all drives (default: false).")
		fmt.Println("      -disk: Specify a disk to search.")
		fmt.Println("             Windows: 'C:\\' or 'D:\\'")
		fmt.Println("             Linux: '/' or '/home/user/'")
		fmt.Println()
		fmt.Println("  find       Find files matching a pattern")
		fmt.Println("    Flags:")
		fmt.Println("      -pattern: File pattern to find (e.g., '*.go').")
		fmt.Println("      -all: Search all drives (default: false).")
		fmt.Println("      -disk: Specify a disk to search.")
		fmt.Println("             Windows: 'C:\\' or 'D:\\'")
		fmt.Println("             Linux: '/' or '/home/user/'")
		fmt.Println()
		fmt.Println("  content    Find files containing specific content")
		fmt.Println("    Flags:")
		fmt.Println("      -string: String to search for in files (e.g., 'TODO').")
		fmt.Println("      -type: File type to search (e.g., '.go').")
		fmt.Println("      -maxsize: Max file size in KB (default: 1024).")
		fmt.Println("      -all: Search all drives (default: false).")
		fmt.Println("      -disk: Specify a disk to search.")
		fmt.Println("             Windows: 'C:\\' or 'D:\\'")
		fmt.Println("             Linux: '/' or '/home/user/'")
		fmt.Println()
		fmt.Println("  extension  Change file extensions")
		fmt.Println("    Flags:")
		fmt.Println("      -pattern: File pattern to change extension (e.g., '*.txt').")
		fmt.Println("      -new: New file extension (e.g., '.md').")
		fmt.Println("      -all: Search all drives (default: false).")
		fmt.Println("      -disk: Specify a disk to search.")
		fmt.Println("             Windows: 'C:\\' or 'D:\\'")
		fmt.Println("             Linux: '/' or '/home/user/'")
		fmt.Println()
		fmt.Println("Use 'file-manager <command> -h' for more information about a command.")
		os.Exit(0)
	}

	// Define subcommands
	countCmd := flag.NewFlagSet("count", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	contentCmd := flag.NewFlagSet("content", flag.ExitOnError)
	extensionCmd := flag.NewFlagSet("extension", flag.ExitOnError)

	// Flags for count
	countPattern := countCmd.String("pattern", "*", "File pattern to count")
	countAll := countCmd.Bool("all", false, "Search all drives")
	countDisk := countCmd.String("disk", "", "Specific disk to search")

	// Flags for remove
	removePattern := removeCmd.String("pattern", "*", "File pattern to remove")
	removeAll := removeCmd.Bool("all", false, "Search all drives")
	removeDisk := removeCmd.String("disk", "", "Specific disk to search")

	// Flags for find
	findPattern := findCmd.String("pattern", "*", "File pattern to find")
	findAll := findCmd.Bool("all", false, "Search all drives")
	findDisk := findCmd.String("disk", "", "Specific disk to search")

	// Flags for content
	contentString := contentCmd.String("string", "", "String to search for in files")
	contentType := contentCmd.String("type", "", "File type to search")
	contentMaxSize := contentCmd.Int("maxsize", 1024, "Max file size in KB")
	contentAll := contentCmd.Bool("all", false, "Search all drives")
	contentDisk := contentCmd.String("disk", "", "Specific disk to search")

	// Flags for extension
	extensionPattern := extensionCmd.String("pattern", "*", "File pattern to change extension")
	newExtension := extensionCmd.String("new", ".txt", "New file extension")
	extensionAll := extensionCmd.Bool("all", false, "Search all drives")
	extensionDisk := extensionCmd.String("disk", "", "Specific disk to search")

	// Parse subcommands
	if len(os.Args) < 2 {
		fmt.Println("File Manager CLI Application")
		fmt.Println("Usage:")
		fmt.Println("  file-manager <command> [flags]")
		fmt.Println("Commands:")
		fmt.Println("  count      Count files matching a pattern")
		fmt.Println("  remove     Remove files matching a pattern")
		fmt.Println("  find       Find files matching a pattern")
		fmt.Println("  content    Find files containing specific content")
		fmt.Println("  extension  Change file extensions")
		fmt.Println("Use 'file-manager <command> -h' for more information about a command.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "count":
		countCmd.Usage = func() {
			fmt.Println("Usage: file-manager count [flags]")
			fmt.Println("Flags:")
			countCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file-manager count -pattern=\"*.txt\" -all")
		}
		countCmd.Parse(os.Args[2:])
		if *countPattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		count, err := file_manipulation.GetFilesCount(*countPattern, *countAll, *countDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Found %d files matching pattern '%s'\n", count, *countPattern)

	case "remove":
		removeCmd.Usage = func() {
			fmt.Println("Usage: file-manager remove [flags]")
			fmt.Println("Flags:")
			removeCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file-manager remove -pattern=\"*.log\" -disk=\"C:\\\"")
		}
		removeCmd.Parse(os.Args[2:])
		if *removePattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		if err := file_manipulation.RemoveFiles(*removePattern, *removeAll, *removeDisk); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Files removed successfully.")

	case "find":
		findCmd.Usage = func() {
			fmt.Println("Usage: file-manager find [flags]")
			fmt.Println("Flags:")
			findCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file-manager find -pattern=\"*.go\" -disk=\"/\"")
		}
		findCmd.Parse(os.Args[2:])
		if *findPattern == "" {
			fmt.Println("Error: File pattern cannot be empty.")
			os.Exit(1)
		}
		files, err := file_manipulation.FindFiles(*findPattern, *findAll, *findDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Found files:", files)

	case "content":
		contentCmd.Usage = func() {
			fmt.Println("Usage: file-manager content [flags]")
			fmt.Println("Flags:")
			contentCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file-manager content -string=\"TODO\" -type=\".go\" -maxsize=512 -all")
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
		files, err := file_manipulation.FindFilesByContent(*contentString, *contentType, *contentMaxSize, *contentAll, *contentDisk)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Found files containing string:", files)

	case "extension":
		extensionCmd.Usage = func() {
			fmt.Println("Usage: file-manager extension [flags]")
			fmt.Println("Flags:")
			extensionCmd.PrintDefaults()
			fmt.Println("Example:")
			fmt.Println("  file-manager extension -pattern=\"*.txt\" -new=\".md\" -disk=\"/\"")
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
		if err := file_manipulation.SetFilesExtension(*extensionPattern, *newExtension, *extensionAll, *extensionDisk); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("File extensions updated successfully.")

	default:
		fmt.Println("Unknown command. Use 'file-manager -h' for help.")
		os.Exit(1)
	}
}
