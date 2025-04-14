package random_utilities

import (
	"strings"
)

// MatchesWildcard checks if a string matches a wildcard pattern.
//
// Description:
// - Supports patterns with "*" as a wildcard (e.g., "test*", "*test", or "*test*").
// - Returns true if the string matches the pattern.
//
// Parameters:
// - pattern (string): The wildcard pattern to match.
// - str (string): The string to check against the pattern.
//
// Returns:
// - bool: True if the string matches the pattern, false otherwise.
//
// Example Usage:
// ```go
// matched := MatchesWildcard("test*", "testing")
// fmt.Println("Matched:", matched) // Output: Matched: true
// ```
func MatchesWildcard(pattern, str string) bool {
	matched := strings.HasPrefix(pattern, "*") || strings.HasSuffix(pattern, "*")
	if matched {
		return strings.Contains(str, strings.Trim(pattern, "*"))
	}
	return pattern == str
}
