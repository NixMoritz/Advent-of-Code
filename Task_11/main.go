package main

import (
	"fmt"
	"strconv"
)

func transformStone(stone int) []int {
	// Rule 1: If stone is 0, replace with 1
	if stone == 0 {
		return []int{1}
	}

	// Rule 2: If even number of digits, split the stone
	stoneStr := strconv.Itoa(stone)
	if len(stoneStr)%2 == 0 {
		mid := len(stoneStr) / 2
		leftNum, _ := strconv.Atoi(stoneStr[:mid])
		rightNum, _ := strconv.Atoi(stoneStr[mid:])
		return []int{leftNum, rightNum}
	}

	// Rule 3: Default - multiply by 2024
	return []int{stone * 2024}
}

func countStonesAfterNBlinks(initialStones []int, nBlinks int) int {
	stoneCounts := make(map[int]int)
	for _, stone := range initialStones {
		stoneCounts[stone]++
	}

	memo := make(map[int][]int)

	for i := 0; i < nBlinks; i++ {
		newStoneCounts := make(map[int]int)
		for stone, count := range stoneCounts {
			var transformed []int
			if res, found := memo[stone]; found {
				transformed = res
			} else {
				transformed = transformStone(stone)
				memo[stone] = transformed
			}

			for _, tStone := range transformed {
				newStoneCounts[tStone] += count
			}
		}
		stoneCounts = newStoneCounts
	}

	totalStones := 0
	for _, count := range stoneCounts {
		totalStones += count
	}

	return totalStones
}

func main() {
	initialStones := []int{9694820, 93, 54276, 1304, 314, 664481, 0, 4}

	nBlinks := 25
	SolvePart1 := countStonesAfterNBlinks(initialStones, nBlinks)

	fmt.Printf("Number of stones after %d blinks: %d\n", nBlinks, SolvePart1)

	nBlinks = 75
	finalStoneCount := countStonesAfterNBlinks(initialStones, nBlinks)

	fmt.Printf("Number of stones after %d blinks: %d\n", nBlinks, finalStoneCount)
}
