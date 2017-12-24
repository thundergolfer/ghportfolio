package main

import (
	"fmt"
	"time"
)

const (
	timelineLimit int = 30
)

func timelineHeader() string {
	curr := time.Now()
	daysInPast := 0
	header := ""

	for daysInPast < timelineLimit {
		header += "|" + curr.Weekday().String()[0:2]
		curr = curr.AddDate(0, 0, -1)
		daysInPast++
	}

	return header
}

// TimelineCount takes a map
func TimelineCount(counts map[string]int) string {
	curr := time.Now()
	daysInPast := 0
	result := ""

	for daysInPast < timelineLimit {
		currKey := timeToDateStr(curr)

		if val, ok := counts[currKey]; ok {
			result += fmt.Sprintf("|%2d", val)
		} else {
			result += "|  "
		}

		curr = curr.AddDate(0, 0, -1)
		daysInPast++
	}

	return result
}
