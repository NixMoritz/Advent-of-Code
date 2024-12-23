package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Part 1: Check if a design can be made using available patterns
func canMakeDesign(design string, patterns []string, memo map[string]bool) bool {
	if design == "" {
		return true
	}

	if result, exists := memo[design]; exists {
		return result
	}

	for _, pattern := range patterns {
		if len(pattern) <= len(design) && design[:len(pattern)] == pattern {
			if canMakeDesign(design[len(pattern):], patterns, memo) {
				memo[design] = true
				return true
			}
		}
	}

	memo[design] = false
	return false
}

// Part 2: Count all possible ways to make a design
func countWaysToMakeDesign(design string, patterns []string, memo map[string]int64) int64 {
	if design == "" {
		return 1
	}

	if count, exists := memo[design]; exists {
		return count
	}

	var totalWays int64 = 0
	for _, pattern := range patterns {
		if len(pattern) <= len(design) && design[:len(pattern)] == pattern {
			totalWays += countWaysToMakeDesign(design[len(pattern):], patterns, memo)
		}
	}

	memo[design] = totalWays
	return totalWays
}

func solve(patterns []string, designs []string) (int, int64) {
	// Part 1: Count possible designs
	possibleCount := 0
	memoPart1 := make(map[string]bool)

	for _, design := range designs {
		if canMakeDesign(design, patterns, memoPart1) {
			possibleCount++
		}
	}

	// Part 2: Count total ways for each possible design
	var totalWays int64 = 0
	memoPart2 := make(map[string]int64)

	for _, design := range designs {
		ways := countWaysToMakeDesign(design, patterns, memoPart2)
		totalWays += ways
	}

	return possibleCount, totalWays
}

func main() {
	// Read input file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read patterns from first line
	var patterns []string
	if scanner.Scan() {
		patternsLine := scanner.Text()
		// Split on comma and space, trim whitespace
		for _, pattern := range strings.Split(patternsLine, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	// Skip empty line
	scanner.Scan()

	// Read designs
	var designs []string
	for scanner.Scan() {
		design := strings.TrimSpace(scanner.Text())
		if design != "" {
			designs = append(designs, design)
		}
	}

	possibleDesigns, totalWays := solve(patterns, designs)
	fmt.Printf("Part 1 - Number of possible designs: %d\n", possibleDesigns)
	fmt.Printf("Part 2 - Total number of different ways: %d\n", totalWays)
}
