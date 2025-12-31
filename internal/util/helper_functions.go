package util

import (
	"time"
)

// Helper to convert category array to string
func CategoryToString(categories []string) string {
	if len(categories) == 0 {
		return ""
	}
	return categories[len(categories)-1]
}

// Helper to parse date string
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
