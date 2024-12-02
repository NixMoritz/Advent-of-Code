package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	amount, err := solveFirst("input_data.txt")
	if err != nil {
		fmt.Errorf("Error solving 2.1: ", err)
	}
	fmt.Printf("Ergebnis Task 2.1 = %v\n", amount)

	amount, err = solveSecond("input_data.txt")
	if err != nil {
		fmt.Errorf("Error solving 2.2: ", err)
	}
	fmt.Printf("Ergebnis Task 2.2 = %v\n", amount)
}

func solveFirst(filename string) (int, error) {
	sol := 0

	data, err := os.ReadFile(filename)
	if err != nil {
		return -1, fmt.Errorf("error reading the file: %w", err)
	}

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(string(line))
		report := make([]int, len(parts))
		for i, part := range parts {
			report[i], err = strconv.Atoi(part)
			if err != nil {
				return 0, fmt.Errorf("error converting to integer: %w", err)
			}
		}

		// Check if the report is safe
		if isSafe(report) {
			sol++
		}
	}

	return sol, nil
}

func isSafe(report []int) bool {
	if len(report) < 2 {
		return false
	}

	isIncreasing := true
	isDecreasing := true

	for i := 0; i < len(report)-1; i++ {
		diff := report[i+1] - report[i]

		if abs(diff) < 1 || abs(diff) > 3 {
			return false
		}
		if diff < 0 {
			isIncreasing = false
		}
		if diff > 0 {
			isDecreasing = false
		}
	}

	return isIncreasing || isDecreasing
}

func abs(absolute int) int {
	if absolute < 0 {
		return -absolute
	}
	return absolute
}

func solveSecond(filename string) (int, error) {
	sol := 0

	data, err := os.ReadFile(filename)
	if err != nil {
		return -1, fmt.Errorf("error reading the file: %w", err)
	}

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(string(line))
		report := make([]int, len(parts))
		for i, part := range parts {
			report[i], err = strconv.Atoi(part)
			if err != nil {
				return 0, fmt.Errorf("error converting to integer: %w", err)
			}
		}

		for i := 0; i < len(report); i++ {
			if isSafe(removeAtIndex(report, i)) {
				sol++
				break
			}
		}
	}

	return sol, nil
}

func removeAtIndex(s []int, index int) []int {
	copy := make([]int, 0, len(s)-1)

	for i := 0; i < len(s); i++ {
		if i != index {
			copy = append(copy, s[i])
		}
	}
	return copy
}
