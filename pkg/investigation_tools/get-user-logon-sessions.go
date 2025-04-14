package investigation_tools

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
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

		// Parse event data
		var eventXML struct {
			EventData struct {
				Data []struct {
					Name  string `xml:"Name,attr"`
					Value string `xml:",chardata"`
				} `xml:"Data"`
			} `xml:"EventData"`
		}

		if err := xml.Unmarshal([]byte(e.StringInserts[0]), &eventXML); err != nil {
			log.Printf("Failed to parse event XML: %v", err)
			continue
		}

		var extractedUserName, extractedLogonType string
		for _, data := range eventXML.EventData.Data {
			switch strings.ToLower(data.Name) {
			case "targetusername":
				extractedUserName = data.Value
			case "logontype":
				extractedLogonType = data.Value
			}
		}

		// Validate extracted data
		if extractedUserName == "" || extractedLogonType == "" {
			log.Println("Incomplete logon event data, skipping...")
			continue
		}

		// Check if event time is within the specified range
		eventTime := e.TimeGenerated
		if eventTime.Before(startDate) || eventTime.After(endDate) {
			continue
		}

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
