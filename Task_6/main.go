package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	/*
		fmt.Println("Test SolvePart1 Result:")
		testSolvePart1()
		result := SolvePart2()
		fmt.Println("SolvePart2 Result:", result)
		fmt.Println("Test SolvePart2 Result:")
		testSolvePart2()
	*/

	data, _ := os.ReadFile("input.txt")
	input := strings.TrimSpace(string(data))

	fmt.Println("SolvePart1 Result:", SolvePart1())
	//fmt.Println("SolvePart2 Result:", SolvePart2())

	fmt.Println(p1(input))
	fmt.Println(p2(input))
}
func testSolvePart1() {
	grid := [][]string{
		{".", ".", ".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", "X", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "X", ".", "."},
		{".", ".", "X", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", ".", "X", "."},
		{"X", ".", "^", ".", ".", ".", ".", ".", ".", "X"},
	}

	sol := guardPathing(grid)
	fmt.Println(sol)
}

func testSolvePart2() {
	grid := [][]string{
		{".", ".", ".", ".", ".", ".", ".", ".", ".", "."},
		{".", "#", "#", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", "#", ".", ".", "#", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", "#", ".", ".", ".", ".", "#", "."},
		{"#", ".", "^", ".", ".", ".", ".", ".", ".", "#"},
	}

	validPositions := findValidObstructionPositions(grid)
	fmt.Println(len(validPositions))

}

func SolvePart1() int {
	grid, err := inputToGridOfChars()
	if err != nil {
		fmt.Println("Error:", err) // Print the error if one occurs
		return -1
	}

	sol := guardPathing(grid)
	fmt.Println("Unique steps before exiting:", sol)
	return sol
}

func SolvePart2() int {
	grid, err := inputToGridOfChars()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	// Find the possible valid obstruction positions
	validPositions := findValidObstructionPositions(grid)
	return len(validPositions)
}

func findValidObstructionPositions(grid [][]string) []string {
	// Find the guard's starting position
	xStart, yStart := findGuard(grid)
	if xStart == -1 || yStart == -1 {
		return nil
	}

	validPositions := []string{}
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == "." && (i != xStart || j != yStart) {
				// Temporarily add the obstruction and simulate the guard's movement
				grid[i][j] = "O"
				fmt.Println("are we bruting? x: ", i, " y: ", j)
				if isGuardStuck(grid, xStart, yStart) {
					validPositions = append(validPositions, fmt.Sprintf("%d,%d", i, j))
				}
				// Remove the obstruction after the test
				grid[i][j] = "."
			}
		}
	}

	return validPositions
}

func isGuardStuck(grid [][]string, xStart, yStart int) bool {
	// Directions: 0 = up, 1 = right, 2 = down, 3 = left
	direction := 0
	visited := make(map[string]bool) // Map to store visited positions
	oCounter := make(map[string]int) // Map to count "O" encounters at positions
	rotationCount := 0               // To track number of rotations without progress

	xPos, yPos := xStart, yStart
	visited[fmt.Sprintf("%d,%d,%d", xPos, yPos, direction)] = true // Include direction in visited positions

	// Keep looping until the guard exits the grid boundaries or gets stuck
	for {
		var nextX, nextY int
		switch direction {
		case 0: // Up
			nextX, nextY = xPos-1, yPos
		case 1: // Right
			nextX, nextY = xPos, yPos+1
		case 2: // Down
			nextX, nextY = xPos+1, yPos
		case 3: // Left
			nextX, nextY = xPos, yPos-1
		}

		// Check if the next position is out of bounds
		if nextX < 0 || nextX >= len(grid) || nextY < 0 || nextY >= len(grid[0]) {
			return false // Exit if out of bounds
		}

		// If the next position is an obstruction, increment the obstruction counter
		if grid[nextX][nextY] == "O" {
			oCounter[fmt.Sprintf("%d,%d", nextX, nextY)]++
		}

		// Check if the guard revisits the same position with the same direction
		if visited[fmt.Sprintf("%d,%d,%d", nextX, nextY, direction)] {
			// If the guard revisits a position with the same direction and has encountered "O" more than once
			if oCounter[fmt.Sprintf("%d,%d", nextX, nextY)] > 1 {
				fmt.Println("Guard stuck at", nextX, nextY) // Debug print
				return true                                 // The guard is stuck in a loop
			}
		}

		// Mark the position as visited with the current direction
		visited[fmt.Sprintf("%d,%d,%d", nextX, nextY, direction)] = true

		// If the next position is a wall or obstruction, the guard needs to turn
		if grid[nextX][nextY] == "#" || grid[nextX][nextY] == "O" {
			// Undo the last movement
			nextX, nextY = xPos, yPos
			// Turn the guard 90° to the right
			direction = (direction + 1) % 4
			rotationCount++

			// If the guard has rotated 4 times without making any progress, break out
			if rotationCount >= 4 {
				fmt.Println("Guard is rotating endlessly, exiting loop.")
				return true
			}
		} else {
			// Reset rotation count if the guard makes a move
			rotationCount = 0
		}

		// Move the guard to the next valid position
		xPos, yPos = nextX, nextY

	}
}

func findGuard(grid [][]string) (int, int) {
	// Search the grid for the guard (denoted by "^")
	for i, line := range grid {
		for j, pos := range line {
			if pos == "^" {
				return i, j // Return the position of the guard
			}
		}
	}
	return -1, -1 // If no guard found, return invalid position
}

func guardPathing(grid [][]string) int {
	// Directions: 0 = up, 1 = right, 2 = down, 3 = left
	direction := 0
	visited := make(map[string]bool) // Map to store visited positions
	xStart, yStart := findGuard(grid)
	if xStart == -1 || yStart == -1 {
		fmt.Println("Guard not found!")
		return -1
	}

	xPos, yPos := xStart, yStart

	// Keep looping until the guard exits the grid boundaries
	for {
		// Check if the position is out of bounds
		if xPos < 0 || xPos >= len(grid) || yPos < 0 || yPos >= len(grid[0]) {
			break // Exit if out of bounds
		}

		// Mark the current position as visited
		//println(xPos, yPos)
		posKey := fmt.Sprintf("%d,%d", xPos, yPos)
		if !visited[posKey] {
			visited[posKey] = true
		}

		// Move the guard based on the current direction
		var nextX, nextY int
		switch direction {
		case 0: // Up
			nextX, nextY = xPos-1, yPos
		case 1: // Right
			nextX, nextY = xPos, yPos+1
		case 2: // Down
			nextX, nextY = xPos+1, yPos
		case 3: // Left
			nextX, nextY = xPos, yPos-1
		}

		// Check if the next position is out of bounds or contains 'X'
		if nextX < 0 || nextX >= len(grid) || nextY < 0 || nextY >= len(grid[0]) {
			return len(visited)

		} else if grid[nextX][nextY] == "#" {
			// Undo the last movement
			nextX, nextY = xPos, yPos

			// Turn the guard 90° to the right
			direction = (direction + 1) % 4
		}

		// Move the guard to the next position
		xPos, yPos = nextX, nextY

	}
	//fmt.Println(visited)
	return len(visited)
}

func inputToGridOfChars() ([][]string, error) {
	// Read the grid from a file
	data, err := os.ReadFile("input.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	grid := make([][]string, len(lines))
	for i, line := range lines {
		grid[i] = make([]string, len(line))
		for j, char := range line {
			grid[i][j] = string(char) // Convert each character to a string
		}
	}
	//fmt.Println(grid)
	return grid, nil
}

func p1(input string) int {
	m, start := parse(input)

	_, path := onLoop(m, start, vec{-1, -1})

	return len(path)
}

func p2(input string) int {
	m, start := parse(input)

	_, orig := onLoop(m, start, vec{-1, -1})
	c := make(chan bool)
	for p := range orig {
		if p == start {
			continue
		}
		go func() {
			ok, _ := onLoop(m, start, p)
			c <- ok
		}()
	}

	loops := 0
	for range len(orig) - 1 {
		if <-c {
			loops++
		}
	}

	return loops

}

func parse(input string) ([][]byte, vec) {
	m := readMatrix(input, func(b byte) byte {
		return b
	})

	start := vec{0, 0}
	R, C := len(m), len(m[0])
outer:
	for r := range R {
		for c := range C {
			if m[r][c] == '^' {
				start = vec{r, c}
				break outer
			}

		}
	}

	return m, start
}

type state struct {
	pos vec
	dir vec
}

func onLoop(m [][]byte, start, obstruction vec) (bool, map[vec]bool) {
	seenPt := map[vec]bool{}
	seenState := map[state]bool{}
	curr := start
	dir := vec{-1, 0}
	for {
		if _, ok := seenState[state{curr, dir}]; ok {
			return true, map[vec]bool{}
		}

		seenPt[curr] = true
		seenState[state{curr, dir}] = true

		next := curr.add(dir)
		r, c := next[0], next[1]

		if !(0 <= r && r < len(m) && c >= 0 && c < len(m[0])) {
			return false, seenPt
		}

		if m[r][c] == '#' || next == obstruction {
			dir = dir.rotate(3)
		} else {
			curr = next
		}
	}

}

/*
utils
*/

func readMatrix[T any](s string, transform func(byte) T) [][]T {
	rows := strings.Split(s, "\n")
	matrix := make([][]T, len(rows))

	for i, row := range rows {
		matrix[i] = make([]T, len(row))
		for j := range row {
			matrix[i][j] = transform(row[j])
		}
	}

	return matrix
}

type vec [2]int

func (u vec) add(v vec) vec {
	return vec{u[0] + v[0], u[1] + v[1]}
}

func (u vec) rotate(n int) vec {
	a, b := u[0], u[1]
	for range n % 4 {
		a, b = -b, a
	}
	return vec{a, b}
}
