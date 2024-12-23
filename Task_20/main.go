package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Point2D struct {
	x, y int
}

func (p Point2D) equals(other Point2D) bool {
	return p.x == other.x && p.y == other.y
}

func (p Point2D) adjacent() []Point2D {
	return []Point2D{
		{p.x, p.y - 1},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
		{p.x - 1, p.y},
	}
}

// Delta represents movement options in the maze
var DELTA = [][]int{
	{0, -2},
	{2, 0},
	{0, 2},
	{-2, 0},
}

type QueueItem struct {
	position Point2D
	time     int
}

type Queue struct {
	items []QueueItem
}

func (q *Queue) push(item QueueItem) {
	q.items = append(q.items, item)
}

func (q *Queue) pop() QueueItem {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) isEmpty() bool {
	return len(q.items) == 0
}

func readInput(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func parsePath(input string) (map[int]map[int]int, error) {
	lines := strings.Split(input, "\n")
	grid := make([][]rune, len(lines))
	var start, end Point2D

	// Parse the input grid
	for y, line := range lines {
		grid[y] = make([]rune, len(line))
		for x, char := range line {
			switch char {
			case 'S':
				start = Point2D{x, y}
				grid[y][x] = '.'
			case 'E':
				end = Point2D{x, y}
				grid[y][x] = '.'
			default:
				grid[y][x] = char
			}
		}
	}

	// Initialize BFS queue
	queue := &Queue{}
	queue.push(QueueItem{start, 0})
	tileSteps := make(map[int]map[int]int)

	// Process queue
	for !queue.isEmpty() && !start.equals(end) {
		item := queue.pop()
		pos, time := item.position, item.time

		// Initialize map for y coordinate if it doesn't exist
		if _, exists := tileSteps[pos.y]; !exists {
			tileSteps[pos.y] = make(map[int]int)
		}
		tileSteps[pos.y][pos.x] = time

		// Check adjacent positions
		for _, adj := range pos.adjacent() {
			if adj.y >= 0 && adj.y < len(grid) && adj.x >= 0 && adj.x < len(grid[0]) {
				if grid[adj.y][adj.x] == '.' {
					if _, exists := tileSteps[adj.y]; !exists {
						tileSteps[adj.y] = make(map[int]int)
					}
					if _, visited := tileSteps[adj.y][adj.x]; !visited {
						queue.push(QueueItem{adj, time + 1})
					}
				}
			}
		}
	}

	return tileSteps, nil
}

func solvePart1(input string) (int, error) {
	tileSteps, err := parsePath(input)
	if err != nil {
		return 0, err
	}

	targetSteps := 30
	if len(tileSteps) > 20 {
		targetSteps = 100
	}

	count := 0
	for y, line := range tileSteps {
		for x, tileTime := range line {
			for _, d := range DELTA {
				nX, nY := x+d[0], y+d[1]
				if nextLine, exists := tileSteps[nY]; exists {
					if nextTime, exists := nextLine[nX]; exists && nextTime > tileTime {
						if nextTime-tileTime-2 >= targetSteps {
							count++
						}
					}
				}
			}
		}
	}

	return count, nil
}

func solvePart2(input string) (int, error) {
	tileSteps, err := parsePath(input)
	if err != nil {
		return 0, err
	}

	targetSteps := 30
	if len(tileSteps) > 20 {
		targetSteps = 100
	}

	count := 0
	for y, line := range tileSteps {
		for x, tileTime := range line {
			for tY, tLine := range tileSteps {
				for tX, tTileTime := range tLine {
					if tTileTime >= tileTime {
						continue
					}

					manhattan := abs(tX-x) + abs(tY-y)
					if manhattan <= 20 && tileTime-tTileTime-manhattan >= targetSteps {
						count++
					}
				}
			}
		}
	}

	return count, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	input, err := readInput("input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	part1Result, err := solvePart1(input)
	if err != nil {
		fmt.Printf("Error solving part 1: %v\n", err)
		return
	}
	fmt.Printf("Part 1 Result: %d\n", part1Result)

	part2Result, err := solvePart2(input)
	if err != nil {
		fmt.Printf("Error solving part 2: %v\n", err)
		return
	}
	fmt.Printf("Part 2 Result: %d\n", part2Result)
}
