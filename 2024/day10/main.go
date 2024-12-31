package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Grid [][]int

func parseGrid(filename string) (Grid, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var grid Grid
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := make([]int, 0)
		for _, ch := range scanner.Text() {
			height, _ := strconv.Atoi(string(ch))
			row = append(row, height)
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return grid, nil
}

func isValidTrail(grid Grid, start, end [2]int) bool {
	rows, cols := len(grid), len(grid[0])
	visited := make(map[[2]int]bool)
	queue := [][2]int{start}
	current := start

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current == end {
			return true
		}

		if visited[current] {
			continue
		}
		visited[current] = true

		directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range directions {
			nr := current[0] + dir[0]
			nc := current[1] + dir[1]

			// Check bounds and height constraint
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
				!visited[[2]int{nr, nc}] &&
				grid[nr][nc] == grid[current[0]][current[1]]+1 {
				queue = append(queue, [2]int{nr, nc})
			}
		}
	}

	return false
}

func findTrailheadScore(grid Grid, trailhead [2]int) int {
	rows, cols := len(grid), len(grid[0])
	score := 0

	// Find all 9-height positions
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 9 {
				if isValidTrail(grid, trailhead, [2]int{r, c}) {
					score++
				}
			}
		}
	}

	return score
}

func solvePartOne(grid Grid) int {
	rows, cols := len(grid), len(grid[0])
	totalScore := 0

	// Find trailheads (positions with height 0)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 0 {
				totalScore += findTrailheadScore(grid, [2]int{r, c})
			}
		}
	}

	return totalScore
}

func findDistinctTrails(grid Grid, trailhead [2]int) int {
	rows, cols := len(grid), len(grid[0])
	distinctTrails := 0

	// Find all 9-height positions
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 9 {
				// Check if trail to this 9-height pos is unique
				if isUniqueTrail(grid, trailhead, [2]int{r, c}) {
					distinctTrails++
				}
			}
		}
	}

	return distinctTrails
}

func isUniqueTrail(grid Grid, start, end [2]int) bool {
	rows, cols := len(grid), len(grid[0])

	// BFS to find paths
	type State struct {
		pos  [2]int
		path [][2]int
	}

	visited := make(map[[2]int]bool)
	queue := []State{{start, [][2]int{start}}}
	foundPaths := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == end {
			foundPaths++
			// If more than one unique path found, return false
			if foundPaths > 1 {
				return false
			}
			continue
		}

		if visited[current.pos] {
			continue
		}
		visited[current.pos] = true

		directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range directions {
			nr := current.pos[0] + dir[0]
			nc := current.pos[1] + dir[1]

			// Check bounds and height constraint
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
				!visited[[2]int{nr, nc}] &&
				grid[nr][nc] == grid[current.pos[0]][current.pos[1]]+1 {
				queue = append(queue, State{
					pos:  [2]int{nr, nc},
					path: append(current.path, [2]int{nr, nc}),
				})
			}
		}
	}

	// Return true if exactly one unique trail found
	return foundPaths == 1
}

func countDistinctTrails(grid Grid, start [2]int) int {
	rows, cols := len(grid), len(grid[0])
	visited := make(map[[2]int]bool)
	directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	var dfs func(pos [2]int, prevHeight int) int
	dfs = func(pos [2]int, prevHeight int) int {
		r, c := pos[0], pos[1]
		// Out of bounds or invalid move
		if r < 0 || r >= rows || c < 0 || c >= cols || visited[pos] || grid[r][c] != prevHeight+1 {
			return 0
		}
		// Found a valid endpoint
		if grid[r][c] == 9 {
			return 1
		}

		// Mark current position as visited
		visited[pos] = true
		distinctPaths := 0

		// Explore all directions
		for _, dir := range directions {
			next := [2]int{r + dir[0], c + dir[1]}
			distinctPaths += dfs(next, grid[r][c])
		}

		// Backtrack
		visited[pos] = false
		return distinctPaths
	}

	return dfs(start, -1) // Start with height -1 (so the first step is always valid)
}

func solvePartTwo(grid Grid) int {
	rows, cols := len(grid), len(grid[0])
	totalRating := 0

	// Find trailheads (positions with height 0)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 0 {
				rating := countDistinctTrails(grid, [2]int{r, c})
				totalRating += rating
			}
		}
	}

	return totalRating
}

func main() {
	grid, err := parseGrid("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	score := solvePartOne(grid)
	fmt.Println("Total trailhead score part 1:", score)

	rating := solvePartTwo(grid)
	fmt.Println("Total trailhead ratings part 2:", rating)
}
