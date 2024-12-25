package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	debug       bool      = false
	debugFile   string    = ""
	debugWriter io.Writer = os.Stdout
)

func debugPrint(format string, a ...interface{}) {
	if debug {
		fmt.Fprintf(debugWriter, format, a...)

	}
}

func parseHeights(schematic string) []int {
	lines := strings.Split(strings.TrimSpace(schematic), "\n")
	if len(lines) == 0 {
		return nil
	}
	numCols := len(lines[0])
	heights := make([]int, numCols)

	isLock := strings.Count(lines[0], "#") == len(lines[0])
	numRows := len(lines)

	for col := 0; col < numCols; col++ {
		height := 0

		if isLock {
			for row := 0; row < numRows; row++ {
				if lines[row][col] == '#' {
					height++
				} else {
					break
				}
			}
		} else {
			for row := numRows - 1; row >= 0; row-- {
				if lines[row][col] == '#' {
					height++
				} else {
					break
				}
			}
		}
		heights[col] = height
	}
	return heights
}

func doesFit(lockHeights, keyHeights []int, totalHeight int) bool {
	if len(lockHeights) != len(keyHeights) {
		return false
	}

	for i := 0; i < len(lockHeights); i++ {
		if lockHeights[i]+keyHeights[i] > totalHeight {
			return false // Early exit
		}
	}
	return true
}

func countFittingPairs(locks [][]int, keys [][]int, totalHeight int) int {
	count := 0
	for i, lock := range locks {
		debugPrint("Lock %d heights: %v\n", i+1, lock)
		for j, key := range keys {
			if doesFit(lock, key, totalHeight) {
				debugPrint("  Fits with Key %d heights: %v\n", j+1, key)
				count++
			}
		}
	}
	return count
}

func parseInput(filename string) ([][]int, [][]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}

	schematics := strings.Split(strings.TrimSpace(string(data)), "\n\n")

	var locks [][]int
	var keys [][]int

	debugPrint("Total schematics found: %d\n\n", len(schematics))

	for i, schematic := range schematics {
		lines := strings.Split(schematic, "\n")
		if len(lines) == 0 {
			continue
		}
		if strings.Count(lines[0], "#") == len(lines[0]) {
			debugPrint("Found lock #%d\n", len(locks)+1)
			locks = append(locks, parseHeights(schematic))
		} else {
			debugPrint("Found key #%d\n", len(keys)+1)
			keys = append(keys, parseHeights(schematic))
		}
		debugPrint("Schematic #%d:\n%v\n\n", i+1, schematic)

	}
	debugPrint("Total locks: %d, Total keys: %d\n\n", len(locks), len(keys))
	return locks, keys, nil
}

func main() {

	// Set debug true to enable debug output, false to disable
	debug = false
	// Set the debugFile variable to enable the output to a file. Set to empty string to disable debug output to a file.
	debugFile = "debug.txt"

	if debug && debugFile != "" {
		file, err := os.Create(debugFile)
		if err != nil {
			log.Fatalf("Failed to create debug file: %v", err)
		}
		defer file.Close()
		debugWriter = file
	}

	solve("example.txt")
	solve("input.txt")
}

func solve(filename string) {
	locks, keys, err := parseInput(filename)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	fmt.Println("\nChecking lock and key pairs for " + filename + ":")
	result := countFittingPairs(locks, keys, 7)
	fmt.Printf("Number of unique fitting lock/key pairs: %d\n", result)
}
