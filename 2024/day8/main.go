package main

import (
	"bufio"
	"fmt"
	"os"
)

const GridSize = 50

type Point struct {
	X, Y int
}

// check if a point is within the grid boundaries.
func (p Point) IsValid() bool {
	return p.X >= 0 && p.X < GridSize && p.Y >= 0 && p.Y < GridSize
}

// add a vector to the current point, returning a new point.
func (p Point) Add(v Vector) Point {
	return Point{X: p.X + v.DX, Y: p.Y + v.DY}
}

// subtract a vector from the current point, returning a new point.
func (p Point) Subtract(v Vector) Point {
	return Point{X: p.X - v.DX, Y: p.Y - v.DY}
}

type Vector struct {
	DX, DY int
}

// create vector from two points.
func NewVector(from, to Point) Vector {
	return Vector{DX: to.X - from.X, DY: to.Y - from.Y}
}

// scale vector by a given factor.
func (v Vector) Scale(factor int) Vector {
	return Vector{DX: v.DX * factor, DY: v.DY * factor}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	frequencies := parseInput(file)

	solvePartOne(frequencies)
	solvePartTwo(frequencies)
}

func parseInput(file *os.File) map[byte][]Point {
	scanner := bufio.NewScanner(file)
	frequencies := make(map[byte][]Point)

	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		for col, char := range line {
			if char != '.' {
				frequencies[byte(char)] = append(frequencies[byte(char)], Point{X: row, Y: col})
			}
		}
	}

	return frequencies
}

func solvePartOne(frequencies map[byte][]Point) {
	antiNodes := make(map[Point]struct{})

	for _, locations := range frequencies {
		for i := 0; i < len(locations)-1; i++ {
			for j := i + 1; j < len(locations); j++ {
				vector := NewVector(locations[i], locations[j])

				// Compute anti-nodes for this pair of points.
				antiNode1 := locations[i].Subtract(vector)
				if antiNode1.IsValid() {
					antiNodes[antiNode1] = struct{}{}
				}

				antiNode2 := locations[j].Add(vector)
				if antiNode2.IsValid() {
					antiNodes[antiNode2] = struct{}{}
				}
			}
		}
	}

	fmt.Println("Part One:", len(antiNodes))
}

func solvePartTwo(frequencies map[byte][]Point) {
	antiNodes := make(map[Point]struct{})

	for _, locations := range frequencies {
		for i := 0; i < len(locations)-1; i++ {
			for j := i + 1; j < len(locations); j++ {
				vector := NewVector(locations[i], locations[j])

				outOfBound := 0
				for period := 0; outOfBound < 2; period++ {
					outOfBound = 0

					antiNode1 := locations[i].Subtract(vector.Scale(period))
					if antiNode1.IsValid() {
						antiNodes[antiNode1] = struct{}{}
					} else {
						outOfBound++
					}

					antiNode2 := locations[j].Add(vector.Scale(period))
					if antiNode2.IsValid() {
						antiNodes[antiNode2] = struct{}{}
					} else {
						outOfBound++
					}
				}
			}
		}
	}

	fmt.Println("Part Two:", len(antiNodes))
}
