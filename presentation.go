package main

import (
	"fmt"
	"math"
	"time"

	"github.com/buger/goterm"
	"github.com/joliv/spark"
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

func timelineCountSparkline(counts map[string]int, yAxisName string) string {
	_ = goterm.NewLineChart(100, 8)
	curr := time.Now()
	daysInPast := 0
	var max, min int
	min = int(math.Inf(1))

	data := []float64{}

	for daysInPast < timelineLimit {
		currKey := timeToDateStr(curr)
		if val, ok := counts[currKey]; ok {
			if val > max {
				max = val
			}
			if val < min {
				min = val
			}
			data = append(data, float64(val))
		} else {
			min = 0
			data = append(data, 0, 0, 0) // zero values take less width
		}

		curr = curr.AddDate(0, 0, -1)
		daysInPast++
	}

	yRange := fmt.Sprintf("  min: %d max: %d", min, max)
	return spark.Line(data) + yRange
}
