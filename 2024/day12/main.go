package main

import (
	"bufio"
	"fmt"
	"os"
)

type Region struct {
	plantType rune
	plots     map[complex128]bool
	area      int
	value     int
	price     int
}

func findRegions(grid [][]rune, calculateValue func(region map[complex128]bool, grid [][]rune) int) []Region {
	rows, cols := len(grid), len(grid[0])
	visited := make(map[complex128]bool)
	var regions []Region

	var exploreRegion func(x, y int) Region
	exploreRegion = func(x, y int) Region {
		plantType := grid[y][x]
		todo := []complex128{complex(float64(x), float64(y))}
		region := Region{
			plantType: plantType,
			plots:     make(map[complex128]bool),
			area:      0,
			value:     0,
		}

		for len(todo) > 0 {
			curr := todo[0]
			todo = todo[1:]
			ix, iy := int(real(curr)), int(imag(curr))
			if visited[curr] {
				continue
			}
			visited[curr] = true
			region.plots[curr] = true
			region.area++

			for _, dir := range []complex128{1, -1, complex(0, 1), complex(0, -1)} {
				nx, ny := ix+int(real(dir)), iy+int(imag(dir))
				if nx >= 0 && nx < cols && ny >= 0 && ny < rows &&
					grid[ny][nx] == plantType && !visited[complex(float64(nx), float64(ny))] {
					todo = append(todo, complex(float64(nx), float64(ny)))
				}
			}
		}

		region.value = calculateValue(region.plots, grid)
		region.price = region.value * region.area
		return region
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			curr := complex(float64(x), float64(y))
			if !visited[curr] {
				regions = append(regions, exploreRegion(x, y))
			}
		}
	}
	return regions
}

func calculateRegionPerimeter(region map[complex128]bool, grid [][]rune) int {
	totalPerimeter := 0
	for cell := range region {
		r, c := int(imag(cell)), int(real(cell))
		plantType := grid[r][c]
		totalPerimeter += calculatePerimeter(r, c, grid, plantType)
	}
	return totalPerimeter
}

func calculatePerimeter(r, c int, grid [][]rune, plantType rune) int {
	rows, cols := len(grid), len(grid[0])
	perimeter := 0
	for _, dir := range []struct{ dr, dc int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		nr, nc := r+dir.dr, c+dir.dc
		if nr < 0 || nr >= rows || nc < 0 || nc >= cols || grid[nr][nc] != plantType {
			perimeter++
		}
	}
	return perimeter
}

func calculateSides(region map[complex128]bool, grid [][]rune) int {
	rows, cols := len(grid), len(grid[0])
	sides := 0

	for coord := range region {
		x, y := int(real(coord)), int(imag(coord))
		for _, dir := range []struct {
			dx, dy int
			x1, y1 int
			x2, y2 int
		}{
			{1, 0, 0, -1, 1, -1},   // Right
			{-1, 0, 0, -1, -1, -1}, // Left
			{0, 1, -1, 0, -1, 1},   // Down
			{0, -1, -1, 0, -1, -1}, // Up
		} {
			nx, ny := x+dir.dx, y+dir.dy
			if nx < 0 || nx >= cols || ny < 0 || ny >= rows || !region[complex(float64(nx), float64(ny))] {
				x1, y1 := x+dir.x1, y+dir.y1
				x2, y2 := x+dir.x2, y+dir.y2
				if region[complex(float64(x1), float64(y1))] && !region[complex(float64(x2), float64(y2))] {
					continue
				}
				sides++
			}
		}
	}
	return sides
}

func readGridFromFile(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return grid, nil
}

func main() {
	grid, err := readGridFromFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	debug := false

	fmt.Println("--- Part One: Perimeter Calculation ---")
	regionsPerimeter := findRegions(grid, calculateRegionPerimeter)
	totalPricePerimeter := 0
	for _, region := range regionsPerimeter {
		if debug {
			fmt.Printf("Region: %c, Area: %d, Perimeter: %d, Price: %d\n", region.plantType, region.area, region.value, region.price)
		}
		totalPricePerimeter += region.price
	}
	fmt.Println("Total Fencing Price (Part One):", totalPricePerimeter)

	fmt.Println("\n--- Part Two: Side Counting ---")
	regionsSides := findRegions(grid, calculateSides)
	totalPriceSides := 0
	for _, region := range regionsSides {
		if debug {
			fmt.Printf("Region: %c, Area: %d, Sides: %d, Price: %d\n", region.plantType, region.area, region.value, region.price)
		}
		totalPriceSides += region.price
	}
	fmt.Println("Total Fencing Price (Part Two):", totalPriceSides)
}
