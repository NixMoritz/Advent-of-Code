package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type coor struct {
	x int
	y int
}

func (c1 coor) Add(c2 coor) coor {
	return coor{
		x: c1.x + c2.x,
		y: c1.y + c2.y,
	}
}

// Constants for grid characters
const (
	robotChar = '@'
	wallChar  = '#'
	emptyChar = '.'
	boxChar   = 'O'
	leftBox   = '['
	rightBox  = ']'
)

// Direction map for character movements
var dirs = map[rune]coor{
	'^': {0, -1},
	'>': {1, 0},
	'v': {0, 1},
	'<': {-1, 0},
}

var stepLeft = coor{-1, 0}
var stepRight = coor{1, 0}

func solvePart1(inputFile string) int {
	grid, moves := loadData(inputFile)

	// Find robot
	var robot coor
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			if grid[i][j] == robotChar {
				robot = coor{j, i}
				break
			}
		}
	}

	var newPos coor
	var newBox coor
	var occupant rune
	for _, move := range moves {

		newPos.x = robot.x + dirs[move].x
		newPos.y = robot.y + dirs[move].y

		occupant = grid[newPos.y][newPos.x]

		// Wall
		if occupant == wallChar {
			continue
		}

		// Unoccupied
		if occupant == emptyChar {
			grid[robot.y][robot.x] = emptyChar
			grid[newPos.y][newPos.x] = robotChar
			robot.x = newPos.x
			robot.y = newPos.y
			continue
		}

		newBox = newPos

		// Box
		for occupant == boxChar {
			newBox.x = newBox.x + dirs[move].x
			newBox.y = newBox.y + dirs[move].y
			occupant = grid[newBox.y][newBox.x]
		}

		if occupant == wallChar {
			continue
		}

		// Push
		for newBox != newPos {
			grid[newBox.y][newBox.x] = boxChar
			newBox.x = newBox.x + dirs[move].x*-1
			newBox.y = newBox.y + dirs[move].y*-1
		}

		grid[robot.y][robot.x] = emptyChar
		grid[newPos.y][newPos.x] = robotChar
		robot.x = newPos.x
		robot.y = newPos.y

	}

	sum := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			if grid[i][j] == boxChar {
				sum += i*100 + j
			}
		}
	}

	return sum
}

func solvePart2(inputFile string) int {
	grid, moves := loadDataPart2(inputFile)

	var robot coor
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			if grid[i][j] == robotChar {
				robot = coor{j, i}
				break
			}
		}
	}

	var newPos coor
	var newBox coor
	var occupant rune
	for _, move := range moves {

		newPos.x = robot.x + dirs[move].x
		newPos.y = robot.y + dirs[move].y

		occupant = grid[newPos.y][newPos.x]

		// Immovable object
		if occupant == wallChar {
			continue
		}

		// Unoccupied
		if occupant == emptyChar {
			grid[robot.y][robot.x] = emptyChar
			grid[newPos.y][newPos.x] = robotChar
			robot.x = newPos.x
			robot.y = newPos.y
			continue
		}

		// Move boxes horizontally
		if move == '<' || move == '>' {
			newBox = newPos
			for occupant == leftBox || occupant == rightBox {
				newBox.x = newBox.x + dirs[move].x
				occupant = grid[newBox.y][newBox.x]
			}

			if occupant == wallChar {
				continue
			}

			// Push
			for newBox != newPos {
				grid[newBox.y][newBox.x] = grid[newBox.y][newBox.x+dirs[move].x*-1]
				newBox.x = newBox.x + dirs[move].x*-1
			}

			grid[robot.y][robot.x] = emptyChar
			grid[newPos.y][newPos.x] = robotChar
			robot.x = newPos.x
			robot.y = newPos.y
			continue
		}

		// Move vertically
		if move == 'v' || move == '^' {
			if !canPush(newPos, dirs[move], grid) {
				continue
			}
			grid = pushBox(newPos, dirs[move], grid)

			grid[robot.y][robot.x] = emptyChar
			grid[newPos.y][newPos.x] = robotChar
			robot.x = newPos.x
			robot.y = newPos.y
		}

	}

	sum := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == leftBox {
				sum += i*100 + j
			}
		}
	}

	return sum
}

func canPush(loc coor, vec coor, grid [][]rune) bool {
	switch grid[loc.y][loc.x] {
	case emptyChar:
		return true
	case leftBox:
		return (canPush(loc.Add(vec), vec, grid) &&
			canPush(loc.Add(stepRight).Add(vec), vec, grid))
	case rightBox:
		return (canPush(loc.Add(vec), vec, grid) &&
			canPush(loc.Add(stepLeft).Add(vec), vec, grid))
	default:
		return false
	}
}

func pushBox(loc coor, vec coor, grid [][]rune) [][]rune {
	if grid[loc.y][loc.x] == emptyChar {
		return grid
	}

	grid = pushBox(
		loc.Add(vec),
		vec,
		grid,
	)

	offset := 1
	if grid[loc.y][loc.x] == rightBox {
		offset = -1
	}

	grid = pushBox(
		loc.Add(coor{x: offset, y: 0}).Add(vec),
		vec,
		grid,
	)

	grid[loc.y+vec.y][loc.x+vec.x] = grid[loc.y][loc.x]
	grid[loc.y+vec.y][loc.x+vec.x+offset] = grid[loc.y][loc.x+offset]

	grid[loc.y][loc.x] = emptyChar
	grid[loc.y][loc.x+offset] = emptyChar

	return grid
}

func loadData(inputFile string) ([][]rune, []rune) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	var grid [][]rune
	var moves []rune
	isGrid := true
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			isGrid = false
			continue
		}
		if isGrid {
			grid = append(grid, []rune(line))
		} else {
			moves = append(moves, []rune(line)...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return grid, moves
}

func loadDataPart2(inputFile string) ([][]rune, []rune) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	var grid [][]rune
	var moves []rune
	isGrid := true
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			isGrid = false
			continue
		}
		if isGrid {
			var newRow []rune
			for _, pos := range line {
				switch pos {
				case emptyChar:
					newRow = append(newRow, emptyChar)
					newRow = append(newRow, emptyChar)
				case boxChar:
					newRow = append(newRow, leftBox)
					newRow = append(newRow, rightBox)
				case wallChar:
					newRow = append(newRow, wallChar)
					newRow = append(newRow, wallChar)
				case robotChar:
					newRow = append(newRow, robotChar)
					newRow = append(newRow, emptyChar)
				}
			}
			grid = append(grid, newRow)
		} else {
			moves = append(moves, []rune(line)...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return grid, moves
}

func main() {

	fmt.Println("Part 1 with Example:", solvePart1("example.txt"))
	fmt.Println("Part 2 with Example:", solvePart2("example.txt"))
	fmt.Println("Part 1 with Input:", solvePart1("input.txt"))
	fmt.Println("Part 2 with Input:", solvePart2("input.txt"))
}
