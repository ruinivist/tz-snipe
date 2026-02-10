package core

import (
	"time"
	"tz-snipe/internal/geodata"
)

func GetTZsMatchingTime(db geodata.TimezoneMap, targetTime string) ([]string, error) {
	// another go quirK: for layout it uses a fixed point in time
	// as reference date Mon, 1/2 03:04:05 PM 2006 -0700
	// 1 2 3 4 5 6 7
	if _, err := time.Parse("15:04", targetTime); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	var matches []string

	for tz := range db {
		// Parse offset to location
		zoneTime, err := time.Parse("-0700", tz)
		if err != nil {
			continue
		}

		// Check if the formatted string matches exactly
		// same minute match
		if now.In(zoneTime.Location()).Format("15:04") == targetTime {
			matches = append(matches, tz)
		}
	}
	return matches, nil
}
