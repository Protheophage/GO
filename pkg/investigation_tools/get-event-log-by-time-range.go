// This module is Windows-specific.

package investigation_tools

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc/eventlog"
)

// EventLog represents a Windows Event Log entry.
type EventLog struct {
	EventID uint32
	Source  string
	Time    time.Time
	Message string
}

// GetEventLogByTimeRange retrieves event logs within a specified time range.
//
// Parameters:
// - startTime (time.Time): The start of the time range.
// - endTime (time.Time): The end of the time range.
// - eventIDs ([]uint32): A list of event IDs to filter (use nil for all IDs).
// - sources ([]string): A list of sources to filter (use nil for all sources).
//
// Returns:
// - []EventLog: A slice of event log entries matching the criteria.
// - error: An error if the event log cannot be read or parsed.
func GetEventLogByTimeRange(startTime, endTime time.Time, eventIDs []uint32, sources []string) ([]EventLog, error) {
	var logs []EventLog

	// Open the Application event log
	elog, err := eventlog.Open("Application")
	if err != nil {
		return nil, fmt.Errorf("failed to open Application event log: %v", err)
	}
	defer elog.Close()

	// Read event logs
	events, err := elog.Read(eventlog.Forwards)
	if err != nil {
		return nil, fmt.Errorf("failed to read event logs: %v", err)
	}

	for _, e := range events {
		// Filter by time range
		if e.TimeGenerated.Before(startTime) || e.TimeGenerated.After(endTime) {
			continue
		}

		// Filter by event ID
		if eventIDs != nil {
			matchesID := false
			for _, id := range eventIDs {
				if e.EventID == id {
					matchesID = true
					break
				}
			}
			if !matchesID {
				continue
			}
		}

		// Filter by source
		if sources != nil {
			matchesSource := false
			for _, source := range sources {
				if e.SourceName == source {
					matchesSource = true
					break
				}
			}
			if !matchesSource {
				continue
			}
		}

		// Add the event log entry to the result
		logs = append(logs, EventLog{
			EventID: e.EventID,
			Source:  e.SourceName,
			Time:    e.TimeGenerated,
			Message: e.StringInserts[0],
		})
	}

	return logs, nil
}
