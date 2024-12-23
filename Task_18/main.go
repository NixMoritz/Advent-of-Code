package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type PriorityQueue []*PathNode

type PathNode struct {
	point    Point
	steps    int
	priority int
	index    int
}

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*PathNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

func readInputFile(filename string) []Point {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		if len(coords) == 2 {
			x, err1 := strconv.Atoi(coords[0])
			y, err2 := strconv.Atoi(coords[1])
			if err1 == nil && err2 == nil {
				points = append(points, Point{x, y})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return points
}

func findShortestPath(corruptedBytes map[Point]bool, gridSize int) int {
	start := Point{0, 0}
	end := Point{gridSize, gridSize}

	// Directions: right, down, left, up
	directions := []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	startNode := &PathNode{
		point:    start,
		steps:    0,
		priority: manhattanDistance(start, end),
	}
	heap.Push(&pq, startNode)

	visited := make(map[Point]int)
	visited[start] = 0

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PathNode)

		if current.point == end {
			return current.steps
		}

		for _, dir := range directions {
			next := Point{current.point.x + dir.x, current.point.y + dir.y}

			// Check grid boundaries
			if next.x < 0 || next.x > gridSize || next.y < 0 || next.y > gridSize {
				continue
			}

			// Check if memory is corrupted
			if corruptedBytes[next] {
				continue
			}

			newSteps := current.steps + 1
			if prevSteps, exists := visited[next]; !exists || newSteps < prevSteps {
				visited[next] = newSteps
				priority := newSteps + manhattanDistance(next, end)
				nextNode := &PathNode{
					point:    next,
					steps:    newSteps,
					priority: priority,
				}
				heap.Push(&pq, nextNode)
			}
		}
	}

	return -1 // No path found
}

func isPathToExitAvailable(corruptedBytes map[Point]bool, gridSize int) bool {
	start := Point{0, 0}
	end := Point{gridSize, gridSize}

	// Directions: right, down, left, up
	directions := []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	startNode := &PathNode{
		point:    start,
		steps:    0,
		priority: manhattanDistance(start, end),
	}
	heap.Push(&pq, startNode)

	visited := make(map[Point]bool)
	visited[start] = true

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PathNode)

		if current.point == end {
			return true
		}

		for _, dir := range directions {
			next := Point{current.point.x + dir.x, current.point.y + dir.y}

			// Check grid boundaries
			if next.x < 0 || next.x > gridSize || next.y < 0 || next.y > gridSize {
				continue
			}

			// Check if memory is corrupted
			if corruptedBytes[next] {
				continue
			}

			// Check if already visited
			if visited[next] {
				continue
			}

			visited[next] = true
			priority := manhattanDistance(next, end)
			nextNode := &PathNode{
				point:    next,
				steps:    current.steps + 1,
				priority: priority,
			}
			heap.Push(&pq, nextNode)
		}
	}

	return false
}

func manhattanDistance(p1, p2 Point) int {
	return int(math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y)))
}

func solvePartOne(inputBytes []Point) int {
	// Track first kilobyte (1024) bytes falling
	corruptedBytes := make(map[Point]bool)
	for i, p := range inputBytes {
		if i < 1024 {
			corruptedBytes[p] = true
		} else {
			break
		}
	}

	// Solve for full 70x70 grid
	gridSize := 70
	return findShortestPath(corruptedBytes, gridSize)
}

func solvePartTwo(inputBytes []Point) string {
	// Track corrupted bytes
	corruptedBytes := make(map[Point]bool)

	// Solve for full 70x70 grid
	gridSize := 70

	// Find the first byte that blocks the path
	for _, p := range inputBytes {
		corruptedBytes[p] = true
		if !isPathToExitAvailable(corruptedBytes, gridSize) {
			return fmt.Sprintf("%d,%d", p.x, p.y)
		}
	}

	return "No blocking byte found"
}

func main() {
	// Read input from file
	inputBytes := readInputFile("input.txt")

	// Solve Part One
	partOneSolution := solvePartOne(inputBytes)
	fmt.Printf("Part One - Shortest Path: %d\n", partOneSolution)

	// Solve Part Two
	partTwoSolution := solvePartTwo(inputBytes)
	fmt.Printf("Part Two - First Blocking Byte: %s\n", partTwoSolution)
}
