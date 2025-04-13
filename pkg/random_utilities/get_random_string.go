// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"math/rand"
	"time"
)

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
