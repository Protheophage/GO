package investigation_tools

import (
	"fmt"
	"log"
	"time"

	"github.com/Protheophage/GO/pkg/random_utilities"
	"golang.org/x/sys/windows/svc/eventlog"
)

// LogonSession represents a user logon session.
type LogonSession struct {
	UserName  string
	LogonTime time.Time
	LogonType string
}

// GetUserLogonSessions retrieves logon sessions from the Windows Security event log.
//
// Description:
// - Filters logon events (Event ID 4624) based on username and time range.
// - Parses event data to extract logon details.
//
// Parameters:
// - userName (string): The username to filter (use "*" for all users).
// - startDate (time.Time): The start of the time range.
// - endDate (time.Time): The end of the time range.
//
// Returns:
// - []LogonSession: A slice of logon session details.
// - error: An error if the event log cannot be read or parsed.
//
// Example Usage:
// ```go
// sessions, err := GetUserLogonSessions("JohnDoe", time.Now().Add(-24*time.Hour), time.Now())
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Logon Sessions:", sessions)
//	}
//
// ```
func GetUserLogonSessions(userName string, startDate, endDate time.Time) ([]LogonSession, error) {
	logonSessions := []LogonSession{}

	// Open the Security event log
	elog, err := eventlog.Open("Security")
	if err != nil {
		return nil, fmt.Errorf("failed to open Security event log: %v", err)
	}
	defer elog.Close()

	// Read event logs
	events, err := elog.Read(eventlog.Forwards)
	if err != nil {
		return nil, fmt.Errorf("failed to read event logs: %v", err)
	}

	for _, e := range events {
		// Filter by Event ID 4624 (Logon Event)
		if e.EventID != 4624 {
			continue
		}

		// Parse event data (replace with actual parsing logic)
		eventTime := e.TimeGenerated
		if eventTime.Before(startDate) || eventTime.After(endDate) {
			continue
		}

		// Extract username and logon type (replace with actual parsing logic)
		extractedUserName := "parsedUserName"   // Replace with actual parsing logic
		extractedLogonType := "parsedLogonType" // Replace with actual parsing logic

		// Check if username matches the filter
		if userName != "*" && !random_utilities.MatchesWildcard(userName, extractedUserName) {
			continue
		}

		logonSessions = append(logonSessions, LogonSession{
			UserName:  extractedUserName,
			LogonTime: eventTime,
			LogonType: extractedLogonType,
		})
	}

	if len(logonSessions) == 0 {
		log.Println("No logon sessions found for the specified criteria.")
	}

	return logonSessions, nil
}
