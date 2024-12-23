package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	possibleCount := 0
	memoPart1 := make(map[string]bool)
	for _, design := range designs {
		if canMakeDesign(design, patterns, memoPart1) {
			possibleCount++
		}
	}
	var totalWays int64 = 0
	memoPart2 := make(map[string]int64)
	for _, design := range designs {
		totalWays += countWaysToMakeDesign(design, patterns, memoPart2)
	}
	return possibleCount, totalWays
}

func processFile(filePath string) (patterns []string, designs []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		patternsLine := scanner.Text()
		for _, pattern := range strings.Split(patternsLine, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	scanner.Scan()
	for scanner.Scan() {
		design := strings.TrimSpace(scanner.Text())
		if design != "" {
			designs = append(designs, design)
		}
	}

	return patterns, designs, nil
}

func main() {
	inputPatterns, inputDesigns, err := processFile("input.txt")
	if err != nil {
		fmt.Println("Error reading input.txt:", err)
		return
	}

	examplePatterns, exampleDesigns, err := processFile("example.txt")
	if err != nil {
		fmt.Println("Error reading example.txt:", err)
		return
	}

	possibleDesignsInput, totalWaysInput := solve(inputPatterns, inputDesigns)
	fmt.Printf("Input.txt - Part 1: %d, Part 2: %d\n", possibleDesignsInput, totalWaysInput)

	possibleDesignsExample, totalWaysExample := solve(examplePatterns, exampleDesigns)
	fmt.Printf("Example.txt - Part 1: %d, Part 2: %d\n", possibleDesignsExample, totalWaysExample)
}
