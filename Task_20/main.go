package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var directions = []Point2D{{0, -2}, {2, 0}, {0, 2}, {-2, 0}}

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

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func parseGrid(input string) ([][]rune, Point2D, Point2D) {
	lines := strings.Split(input, "\n")
	grid := make([][]rune, len(lines))
	var start, end Point2D

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

	return grid, start, end
}

func bfs(grid [][]rune, start, end Point2D) map[Point2D]int {
	queue := &Queue{}
	queue.push(QueueItem{start, 0})
	tileSteps := make(map[Point2D]int)

	for !queue.isEmpty() {
		item := queue.pop()
		pos, time := item.position, item.time

		if _, visited := tileSteps[pos]; visited {
			continue
		}
		tileSteps[pos] = time

		for _, adj := range pos.adjacent() {
			if adj.y >= 0 && adj.y < len(grid) && adj.x >= 0 && adj.x < len(grid[0]) && grid[adj.y][adj.x] == '.' {
				queue.push(QueueItem{adj, time + 1})
			}
		}
	}

	return tileSteps
}

func solvePart1(input string) (int, error) {
	grid, start, end := parseGrid(input)
	tileSteps := bfs(grid, start, end)

	targetSteps := 30
	if len(grid) > 20 {
		targetSteps = 100
	}

	count := 0

	for pos, tileTime := range tileSteps {
		for _, d := range directions {
			nextPos := Point2D{pos.x + d.x, pos.y + d.y}
			if nextTime, exists := tileSteps[nextPos]; exists && nextTime > tileTime {
				if nextTime-tileTime-2 >= targetSteps {
					count++
				}
			}
		}
	}

	return count, nil
}

func solvePart2(input string) (int, error) {
	grid, start, end := parseGrid(input)
	tileSteps := bfs(grid, start, end)

	targetSteps := 30
	if len(grid) > 20 {
		targetSteps = 100
	}

	count := 0

	for pos, tileTime := range tileSteps {
		for otherPos, otherTime := range tileSteps {
			if otherTime >= tileTime {
				continue
			}

			manhattan := abs(pos.x-otherPos.x) + abs(pos.y-otherPos.y)
			if manhattan <= 20 && tileTime-otherTime-manhattan >= targetSteps {
				count++
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

	part2Result, err := solvePart2(input)
	if err != nil {
		fmt.Printf("Error solving part 2: %v\n", err)
		return
	}

	fmt.Printf("Input.txt - Part 1: %d, Part 2: %d\n", part1Result, part2Result)
}
