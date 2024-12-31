package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const (
	UP byte = iota
	DOWN
	LEFT
	RIGHT
)

var directionBytes = []byte{UP, DOWN, LEFT, RIGHT}

var directions = [4][2]int{
	UP:    {-1, 0},
	DOWN:  {1, 0},
	LEFT:  {0, -1},
	RIGHT: {0, 1},
}

type Point int64

func newPoint(i, j int) Point {
	return Point((int64(i) << 32) | int64(j))
}

func (p Point) coords() (int, int) {
	return int(int64(p) >> 32), int(int64(p) & 0xFFFFFFFF)
}

type Map struct {
	grid  [][]byte
	start Point
	end   Point
	rows  int
	cols  int
}

type PathPoint struct {
	score     int
	direction byte
	pos       Point
}

type PQ []PathPoint

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].score < pq[j].score }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(PathPoint))
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type Result struct {
	score int
	paths map[Point]struct{}
}

func main() {
	m := readMap("example.txt")
	if m == nil {
		return
	}

	fmt.Println("Part 1 with example:", part1(m))
	fmt.Println("Part 2 with example:", part2(m))

	m = readMap("input.txt")
	if m == nil {
		return
	}
	fmt.Println("Part 1:", part1(m))
	fmt.Println("Part 2:", part2(m))
}

func readMap(filename string) *Map {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	m := &Map{}
	scanner := bufio.NewScanner(file)

	m.grid = make([][]byte, 0, 1000)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]byte, len(line))
		copy(row, line)
		m.grid = append(m.grid, row)
	}

	m.rows = len(m.grid)
	m.cols = len(m.grid[0])

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			switch m.grid[i][j] {
			case 'S':
				m.start = newPoint(i, j)
			case 'E':
				m.end = newPoint(i, j)
			}
		}
	}

	return m
}

func (m *Map) isWall(pos Point) bool {
	i, j := pos.coords()
	return m.grid[i][j] == '#'
}

func findPaths(m *Map) Result {
	pq := make(PQ, 0, m.rows*m.cols)
	heap.Init(&pq)

	previous := make(map[PathPoint][]PathPoint, m.rows*m.cols)
	scores := make(map[PathPoint]int, m.rows*m.cols)

	startPoint := PathPoint{0, RIGHT, m.start}
	scores[startPoint] = 0
	heap.Push(&pq, startPoint)

	shortestScore := -1
	var allPathPoints map[Point]struct{}

	for pq.Len() > 0 {
		pathPoint := heap.Pop(&pq).(PathPoint)

		if pathPoint.pos == m.end {
			if shortestScore == -1 || pathPoint.score == shortestScore {
				shortestScore = pathPoint.score
				allPathPoints = reconstructPaths(previous, pathPoint)
			} else if pathPoint.score > shortestScore {
				break
			}
			continue
		}

		i, j := pathPoint.pos.coords()

		for _, dir := range directionBytes {
			di, dj := directions[dir][0], directions[dir][1]
			newI, newJ := i+di, j+dj

			// Bounds check
			if newI < 0 || newI >= m.rows || newJ < 0 || newJ >= m.cols {
				continue
			}

			newPos := newPoint(newI, newJ)
			if m.isWall(newPos) {
				continue
			}

			if (pathPoint.direction == UP && dir == DOWN) ||
				(pathPoint.direction == DOWN && dir == UP) ||
				(pathPoint.direction == LEFT && dir == RIGHT) ||
				(pathPoint.direction == RIGHT && dir == LEFT) {
				continue
			}

			newScore := pathPoint.score + 1
			if pathPoint.direction != dir {
				newScore += 1000
			}

			newPathPoint := PathPoint{newScore, dir, newPos}
			if prevScore, ok := scores[newPathPoint]; !ok || newScore < prevScore {
				scores[newPathPoint] = newScore
				previous[newPathPoint] = []PathPoint{pathPoint}
				heap.Push(&pq, newPathPoint)
			} else if newScore == prevScore {
				previous[newPathPoint] = append(previous[newPathPoint], pathPoint)
			}
		}
	}

	return Result{score: shortestScore, paths: allPathPoints}
}

func reconstructPaths(previous map[PathPoint][]PathPoint, start PathPoint) map[Point]struct{} {
	path := make(map[Point]struct{}, len(previous))
	queue := make([]PathPoint, 0, len(previous))
	queue = append(queue, start)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		path[curr.pos] = struct{}{}
		queue = append(queue, previous[curr]...)
	}

	return path
}

func part1(m *Map) int {
	return findPaths(m).score
}

func part2(m *Map) int {
	return len(findPaths(m).paths)
}
