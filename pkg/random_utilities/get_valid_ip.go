package random_utilities

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// GetValidIP prompts the user to enter a valid IP address.
//
// Description:
// - Continuously prompts the user until a valid IPv4 address is entered.
// - Validates the input using a regular expression.
//
// Parameters: None
//
// Returns:
// - string: The valid IP address entered by the user.
// - error: An error if the input or validation fails.
//
// Example Usage:
// ```go
// ip, err := GetValidIP()
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Valid IP address:", ip)
//	}
//
// ```
func GetValidIP() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	ipRegex := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

	for {
		fmt.Print("Enter an IP address: ")
		ipAddress, _ := reader.ReadString('\n')
		ipAddress = strings.TrimSpace(ipAddress) // Remove newline and surrounding whitespace

		if ipRegex.MatchString(ipAddress) {
			fmt.Printf("The IP address is %s\n", ipAddress)
			return ipAddress, nil
		} else {
			fmt.Println("Invalid IP address. Please try again.")
		}
	}
}
