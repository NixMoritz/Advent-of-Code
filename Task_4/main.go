package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Solve First Result:", solveFirst())

	fmt.Println("Solve Second Result:", solveSecond())
}

func solveFirst() int {
	grid, err := inputToGridOfChars()
	if err != nil {
		return -1
	}

	word := "XMAS"
	reversedWord := reverseString(word)
	counts := []int{
		countHorizontal(grid, word),
		countHorizontal(grid, reversedWord),
		countVertical(grid, word),
		countVertical(grid, reversedWord),
		countDiagonals(grid, word),
		countDiagonals(grid, reversedWord),
	}

	return sum(counts)
}

func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func solveSecond() int {
	grid, err := inputToGridOfChars()
	if err != nil {
		return -1
	}
	count := 0

	for y, row := range grid {
		for x, char := range row {
			if char == 'M' {
				count += countMASPatterns(grid, x, y)
			}
		}
	}

	return count / 2
}

func inputToGridOfChars() ([][]rune, error) {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid, nil
}

func countMASPatterns(grid [][]rune, x, y int) int {
	patterns := []func([][]rune, int, int) bool{
		checkMASBottomRight,
		checkMASBottomLeft,
		checkMASTopRight,
		checkMASTopLeft,
	}

	count := 0
	for _, pattern := range patterns {
		if pattern(grid, x, y) {
			count++
		}
	}

	return count
}

func checkMASBottomRight(grid [][]rune, x, y int) bool {
	return x+2 < len(grid[0]) && y+2 < len(grid) &&
		grid[y+1][x+1] == 'A' && grid[y+2][x+2] == 'S' &&
		((grid[y][x+2] == 'S' && grid[y+2][x] == 'M') ||
			(grid[y][x+2] == 'M' && grid[y+2][x] == 'S'))
}

func checkMASBottomLeft(grid [][]rune, x, y int) bool {
	return x-2 >= 0 && y+2 < len(grid) &&
		grid[y+1][x-1] == 'A' && grid[y+2][x-2] == 'S' &&
		((grid[y][x-2] == 'S' && grid[y+2][x] == 'M') ||
			(grid[y][x-2] == 'M' && grid[y+2][x] == 'S'))
}

func checkMASTopRight(grid [][]rune, x, y int) bool {
	return x+2 < len(grid[0]) && y-2 >= 0 &&
		grid[y-1][x+1] == 'A' && grid[y-2][x+2] == 'S' &&
		((grid[y][x+2] == 'S' && grid[y-2][x] == 'M') ||
			(grid[y][x+2] == 'M' && grid[y-2][x] == 'S'))
}

func checkMASTopLeft(grid [][]rune, x, y int) bool {
	return x-2 >= 0 && y-2 >= 0 &&
		grid[y-1][x-1] == 'A' && grid[y-2][x-2] == 'S' &&
		((grid[y][x-2] == 'S' && grid[y-2][x] == 'M') ||
			(grid[y][x-2] == 'M' && grid[y-2][x] == 'S'))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func countHorizontal(grid [][]rune, target string) int {
	count := 0
	for _, line := range grid {
		for i := 0; i <= len(line)-len(target); i++ {
			if string(line[i:i+len(target)]) == target {
				count++
			}
		}
	}
	return count
}

func countVertical(grid [][]rune, target string) int {
	count := 0
	length := len(target)
	for col := 0; col < len(grid[0]); col++ {
		for row := 0; row <= len(grid)-length; row++ {
			match := true
			for i := 0; i < length; i++ {
				if grid[row+i][col] != rune(target[i]) {
					match = false
					break
				}
			}
			if match {
				count++
			}
		}
	}
	return count
}

func countDiagonals(grid [][]rune, target string) int {
	count := 0
	length := len(target)

	for row := 0; row <= len(grid)-length; row++ {
		for col := 0; col <= len(grid[0])-length; col++ {
			match := true
			for i := 0; i < length; i++ {
				if grid[row+i][col+i] != rune(target[i]) {
					match = false
					break
				}
			}
			if match {
				count++
			}
		}
	}

	for row := 0; row <= len(grid)-length; row++ {
		for col := length - 1; col < len(grid[0]); col++ {
			match := true
			for i := 0; i < length; i++ {
				if grid[row+i][col-i] != rune(target[i]) {
					match = false
					break
				}
			}
			if match {
				count++
			}
		}
	}

	return count
}
