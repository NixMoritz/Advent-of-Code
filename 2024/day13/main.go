package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Machine struct {
	AX, AY         int
	BX, BY         int
	PrizeX, PrizeY int
}

const OFFSET = 10000000000000

// SolveLinearEquations solves a system of two linear equations for Part 2
func SolveLinearEquations(coefficients []int64) (int64, int64, bool) {
	a1, b1, c1 := coefficients[0], coefficients[2], coefficients[4]
	a2, b2, c2 := coefficients[1], coefficients[3], coefficients[5]

	determinant := a1*b2 - a2*b1
	if determinant == 0 {
		return 0, 0, false // No unique solution
	}

	numeratorA := c1*b2 - c2*b1
	numeratorB := a1*c2 - a2*c1

	if numeratorA%determinant != 0 || numeratorB%determinant != 0 {
		return 0, 0, false // No integer solution
	}

	a := numeratorA / determinant
	b := numeratorB / determinant

	return a, b, true
}

// Part 1 solution using brute force approach
func findPossibleTokenCombination(machine Machine) (int, bool) {
	ax, ay, bx, by := machine.AX, machine.AY, machine.BX, machine.BY
	px, py := machine.PrizeX, machine.PrizeY

	for a := 0; a <= 100; a++ {
		remainingX := px - a*ax
		remainingY := py - a*ay

		if remainingX%bx == 0 && remainingY%by == 0 {
			b := remainingX / bx
			if b >= 0 && b <= 100 && b*by == remainingY {
				tokens := a*3 + b*1
				return tokens, true
			}
		}
	}

	return 0, false
}

func solveClawMachinePrizes(machines []Machine) int {
	totalTokens := 0
	prizesWon := 0

	for _, machine := range machines {
		tokens, canWin := findPossibleTokenCombination(machine)
		if canWin {
			totalTokens += tokens
			prizesWon++
		}
	}

	fmt.Printf("Prizes Won: %d\n", prizesWon)
	return totalTokens
}

func parsePart1(scanner *bufio.Scanner) []Machine {
	var machines []Machine
	reA := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	reB := regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	rePrize := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	var currentMachine Machine
	machineLines := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matchA := reA.FindStringSubmatch(line)
		matchB := reB.FindStringSubmatch(line)
		matchPrize := rePrize.FindStringSubmatch(line)

		if matchA != nil {
			currentMachine.AX, _ = strconv.Atoi(matchA[1])
			currentMachine.AY, _ = strconv.Atoi(matchA[2])
			machineLines++
		} else if matchB != nil {
			currentMachine.BX, _ = strconv.Atoi(matchB[1])
			currentMachine.BY, _ = strconv.Atoi(matchB[2])
			machineLines++
		} else if matchPrize != nil {
			currentMachine.PrizeX, _ = strconv.Atoi(matchPrize[1])
			currentMachine.PrizeY, _ = strconv.Atoi(matchPrize[2])
			machineLines++
		}

		if machineLines == 3 {
			machines = append(machines, currentMachine)
			machineLines = 0
		}
	}

	return machines
}

func parsePart2(scanner *bufio.Scanner) [][]int64 {
	var games [][]int64
	game := []int64{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		var val1, val2 int64
		fmt.Sscanf(parts[len(parts)-2][2:], "%d", &val1)
		fmt.Sscanf(parts[len(parts)-1][2:], "%d", &val2)

		game = append(game, val1, val2)
		if len(game) == 6 {
			game[4] += OFFSET
			game[5] += OFFSET
			games = append(games, game)
			game = []int64{}
		}
	}

	return games
}

func parseInput(filename string, part2 bool) ([]Machine, [][]int64) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if part2 {
		return nil, parsePart2(scanner)
	}
	return parsePart1(scanner), nil
}

func main() {
	// Solve Part 1
	machines, _ := parseInput("input.txt", false)
	result := solveClawMachinePrizes(machines)
	fmt.Println("Part 1 - Minimum tokens to win prizes:", result)

	// Solve Part 2
	_, games := parseInput("input.txt", true)
	total := int64(0)

	for _, game := range games {
		a, b, valid := SolveLinearEquations(game)
		if valid {
			result := a*3 + b
			total += result
		}
	}

	fmt.Println("Part 2 - Minimum tokens with offset:", total)
}
