package main

import (
	"time"
)

// getCurrentYear returns the current year as an integer
func getCurrentYear() int {
	return time.Now().Year()
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// loadNeighborhoods initializes and returns a slice of neighborhoods with their characteristics
func loadNeighborhoods() []Neighborhood {
	return []Neighborhood{
		{"Downtown", 1.5, 6, [2]float64{0.1, 2.0}, [2]float64{5.0, 15.0}, [2]int{85, 100}, [2]float64{200, 600}},
		{"Midtown", 1.2, 7, [2]float64{2.0, 5.0}, [2]float64{3.0, 8.0}, [2]int{70, 90}, [2]float64{150, 400}},
		{"NorthSide", 1.3, 8, [2]float64{5.0, 10.0}, [2]float64{1.0, 4.0}, [2]int{50, 75}, [2]float64{100, 300}},
		{"SouthSide", 0.9, 5, [2]float64{4.0, 9.0}, [2]float64{4.0, 12.0}, [2]int{40, 65}, [2]float64{75, 250}},
		{"Eastside", 1.1, 6, [2]float64{6.0, 12.0}, [2]float64{2.0, 7.0}, [2]int{30, 60}, [2]float64{50, 200}},
		{"Westside", 1.0, 7, [2]float64{7.0, 15.0}, [2]float64{2.0, 6.0}, [2]int{25, 55}, [2]float64{25, 150}},
		{"SuburbsNorth", 1.4, 9, [2]float64{12.0, 20.0}, [2]float64{0.5, 2.0}, [2]int{15, 40}, [2]float64{150, 350}},
		{"SuburbsSouth", 1.25, 8, [2]float64{10.0, 18.0}, [2]float64{1.0, 3.0}, [2]int{10, 35}, [2]float64{125, 300}},
	}
}
