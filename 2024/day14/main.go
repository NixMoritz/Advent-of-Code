package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Robot struct {
	x, y, vx, vy int
}

// Extract integers from input line using regex
func extractIntegers(line string) []int {
	numRegex := regexp.MustCompile(`-?\d+`)
	matches := numRegex.FindAllString(line, -1)
	var numbers []int
	for _, match := range matches {
		num, _ := strconv.Atoi(match)
		numbers = append(numbers, num)
	}
	return numbers
}

// Parse input into a list of robots
func parseInput(input string) []Robot {
	var robots []Robot
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		ints := extractIntegers(line)
		robots = append(robots, Robot{x: ints[0], y: ints[1], vx: ints[2], vy: ints[3]})
	}
	return robots
}

// Simulate robot movement for a given time step
func simulateRobots(robots []Robot, width, height, steps int) []Robot {
	for i := range robots {
		robots[i].x = (robots[i].x + robots[i].vx*steps + width) % width
		robots[i].y = (robots[i].y + robots[i].vy*steps + height) % height
	}
	return robots
}

// Solve Part 1: Calculate quadrant product
func solvePart1(input string) int {
	robots := parseInput(input)
	gridWidth, gridHeight := 101, 103
	if len(strings.Split(input, "\n")) <= 20 {
		gridWidth, gridHeight = 11, 7
	}

	simulateRobots(robots, gridWidth, gridHeight, 100)

	midX, midY := gridWidth/2, gridHeight/2
	quadrants := make([]int, 4)

	for _, robot := range robots {
		if robot.x != midX && robot.y != midY {
			quadrant := 0
			if robot.x > midX {
				quadrant++
			}
			if robot.y > midY {
				quadrant += 2
			}
			quadrants[quadrant]++
		}
	}

	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

// Solve Part 2: Find first second with unique positions
func solvePart2(input string) int {
	robots := parseInput(input)
	gridWidth, gridHeight := 101, 103
	seconds := 0

	for {
		seconds++
		set := make(map[string]bool)
		robots = simulateRobots(robots, gridWidth, gridHeight, 1)

		for _, robot := range robots {
			key := fmt.Sprintf("%d,%d", robot.x, robot.y)
			set[key] = true
		}

		if len(set) == len(robots) {
			return seconds
		}
	}
}

// Save the grid as a JPEG image
func saveGridAsJPEG(robots []Robot, width, height, second int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	for _, robot := range robots {
		img.Set(robot.x, robot.y, color.RGBA{0, 0, 0, 255})
	}

	dir := "output"
	os.MkdirAll(dir, os.ModePerm)
	filePath := filepath.Join(dir, fmt.Sprintf("grid_%04d.jpeg", second))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Printf("Error encoding JPEG: %v\n", err)
		return
	}

	fmt.Printf("Saved grid image for second %d as %s\n", second, filePath)
}

// Generate JPEG grids for Part 2 visualization
func generatePart2JPEG(input string, width, height, limit int) {
	robots := parseInput(input)
	for time := 0; time < limit; time++ {
		positions := simulateRobots(robots, width, height, time)
		saveGridAsJPEG(positions, width, height, time)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Part2 visual mode: Y/N>")
		os.Exit(1)
	}

	inputBytes, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Error reading input.txt:", err)
		return
	}
	input := string(inputBytes)

	part1Result := solvePart1(input)
	fmt.Printf("Part 1 Result: %d\n", part1Result)

	if strings.ToLower(os.Args[1]) == "y" {
		generatePart2JPEG(input, 101, 103, 10000)
	}

	part2Result := solvePart2(input)
	fmt.Printf("Part 2 Result: %d\n", part2Result)
}
