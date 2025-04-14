// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"math/rand"
	"time"
)

// GetRandomString generates a random string of the specified length.
//
// Description:
// - Uses a predefined set of characters to generate the string.
// - Defaults to a length of 14 if the provided length is less than or equal to 0.
//
// Parameters:
// - length (int): The length of the random string to generate.
//
// Returns:
// - string: The generated random string.
//
// Example Usage:
// ```go
// randomString := GetRandomString(10)
// fmt.Println("Random string:", randomString)
// ```
func GetRandomString(length int) string {
	if length <= 0 {
		length = 14
	}

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!\"#$%&'()*+,-./")
	rand.Seed(time.Now().UnixNano())

	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
