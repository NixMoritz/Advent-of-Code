package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type (
	codeSequence struct {
		code        string
		numericPart int
	}
	keypad [][]string
	button struct {
		y     int
		x     int
		value rune
	}
	direction struct {
		yOffset int
		xOffset int
		value   rune
	}
	move struct {
		from      button
		to        button
		direction direction
	}
	point struct {
		y    int
		x    int
		path string
	}
	pathKey struct {
		seq        string
		robotDepth int
	}
	searchKey struct {
		start, end string
	}
)

const (
	EMPTY = ' '
)

var (
	codeSequences  []*codeSequence
	numberPanel    = keypad{{"7", "8", "9"}, {"4", "5", "6"}, {"1", "2", "3"}, {" ", "0", "A"}}
	directionPanel = keypad{{" ", "^", "A"}, {"<", "v", ">"}}
	directionMap   = map[rune]direction{
		'>': {0, 1, '>'}, 'v': {1, 0, 'v'}, '<': {0, -1, '<'}, '^': {-1, 0, '^'},
	}
	instructionCache = make(map[pathKey]int)
	pathsCache       = make(map[searchKey][]string)
)

func (k keypad) findValue(v string) (int, int) {
	for y, row := range k {
		for x, value := range row {
			if value == v {
				return y, x
			}
		}
	}
	return -1, -1
}

func (k keypad) shortestPaths(start, end string) []string {
	startY, startX := k.findValue(start)
	blankY, blankX := k.findValue(" ")
	endY, endX := k.findValue(end)

	key := searchKey{start, end}
	if prev, seen := pathsCache[key]; seen {
		return prev
	}

	paths := []string{}
	queue := []point{{y: startY, x: startX, path: ""}}
	visited := make(map[[2]int]struct{})
	visited[[2]int{startY, startX}] = struct{}{}
	visited[[2]int{blankY, blankX}] = struct{}{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.y == endY && curr.x == endX {
			paths = append(paths, curr.path+"A")
			continue
		}

		for char, direction := range directionMap {
			nextY, nextX := curr.y+direction.yOffset, curr.x+direction.xOffset
			if nextY >= 0 && nextY < len(k) && nextX >= 0 && nextX < len(k[nextY]) {
				if _, seen := visited[[2]int{nextY, nextX}]; !seen {
					queue = append(queue, point{y: nextY, x: nextX, path: curr.path + string(char)})
				}
			}
		}
		visited[[2]int{curr.y, curr.x}] = struct{}{}
	}

	pathsCache[key] = paths
	return paths
}

func init() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var numericVal int
		fmt.Sscanf(line, "%dA", &numericVal)
		codeSequences = append(codeSequences, &codeSequence{code: line, numericPart: numericVal})
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("Execution time: %v\n", time.Since(start))
	}
}

func min(s []int) int {
	minVal := s[0]
	for _, v := range s[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func minLengthOpDirPanel(seq string, numberRobots int) int {
	if numberRobots == 0 {
		return len(seq)
	}

	key := pathKey{seq, numberRobots}
	if cached, known := instructionCache[key]; known {
		return cached
	}

	result := 0
	currentlyAt := "A"
	for _, character := range seq {
		charAsString := string(character)
		paths := directionPanel.shortestPaths(currentlyAt, charAsString)
		possibleOptions := []int{}
		for _, subSeq := range paths {
			possibleOptions = append(possibleOptions, minLengthOpDirPanel(subSeq, numberRobots-1))
		}
		result += min(possibleOptions)
		currentlyAt = charAsString
	}

	instructionCache[key] = result
	return result
}

func (cs *codeSequence) calcSequence(seq string, numberRobots int) int {
	result := 0
	currentlyAt := "A"
	for _, character := range seq {
		charAsString := string(character)
		paths := numberPanel.shortestPaths(currentlyAt, charAsString)
		possibleOptions := []int{}
		for _, subSeq := range paths {
			possibleOptions = append(possibleOptions, minLengthOpDirPanel(subSeq, numberRobots))
		}
		result += min(possibleOptions)
		currentlyAt = charAsString
	}
	return result
}

func partOne() int {
	result := 0
	for _, seq := range codeSequences {
		result += seq.calcSequence(seq.code, 2) * seq.numericPart
	}
	return result
}

func partTwo() int {
	result := 0
	for _, seq := range codeSequences {
		result += seq.calcSequence(seq.code, 25) * seq.numericPart
	}
	return result
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
