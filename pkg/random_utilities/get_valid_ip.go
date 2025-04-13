package random_utilities

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

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
