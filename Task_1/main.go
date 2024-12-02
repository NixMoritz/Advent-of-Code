package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func solve(left []int, right []int) int {

	distance := 0
	left, right = sorting(left), sorting(right)
	if len(left) == len(right) {
		for idx, num := range left {
			distance += abs(num - right[idx])
		}
	}
	return distance
}

func sorting(slice []int) []int {
	sorted := false

	for !sorted {
		sorted = true

		for i := 0; i < len(slice)-1; i++ {
			if slice[i] > slice[i+1] {
				slice[i], slice[i+1] = slice[i+1], slice[i]
				sorted = false
			}
		}
	}

	//fmt.Printf("Ergebnis = %v\n", slice)
	return slice
}

func abs(absolute int) int {
	if absolute < 0 {
		return -absolute
	}
	return absolute
}

func main() {

	left, right, err := fillArrays("taskData.txt")
	if err != nil {
		println("Error occurred")
	}
	sol1 := solve(left, right)
	sol2 := simPointers(left, right)
	sol3 := findSimilarities(left, right)
	sol4 := findSimilaritiesMap(left, right)
	fmt.Printf("Ergebnis Task 1.1 = %v\n", sol1)
	fmt.Printf("Ergebnis Pointers Task 1.2 = %v\n", sol2)
	fmt.Printf("Ergebnis Two Arrays Task 1.2 = %v\n", sol3)
	fmt.Printf("Ergebnis Map Task 1.2 = %v\n", sol4)
}

func fillArrays(filename string) ([]int, []int, error) {
	oneArray, secondArray := []int{}, []int{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading the file: %w", err)
	}

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(string(line))
		if len(parts) < 2 {
			return nil, nil, fmt.Errorf("invalid line format: %s", string(line))
		}

		val1, err1 := strconv.Atoi(parts[0])
		val2, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("invalid number in line: %s", string(line))
		}

		oneArray = append(oneArray, val1)
		secondArray = append(secondArray, val2)
	}

	return oneArray, secondArray, nil
}

func simPointers(left []int, right []int) int {

	sol := 0
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			i++
		}

		if left[i] > right[j] {
			j++

		} else {
			count := 0
			num := left[i]

			for j < len(right) && right[j] == num {
				count++
				j++
			}

			sol += num * count
		}
	}

	return sol
}
func findSimilarities(left []int, right []int) int {
	left, right = sorting(left), sorting(right)
	sol := 0

	// inefficient af
	for _, num := range left {
		count := 0
		for _, num2 := range right {
			if num == num2 {
				count++
			}
		}
		sol += num * count
	}
	return sol
}

func findSimilaritiesMap(left []int, right []int) int {
	counts := make(map[int]int)
	for _, num := range right {
		counts[num]++
	}

	sol := 0
	for _, num := range left {
		if count, exists := counts[num]; exists {
			sol += num * count
		}
	}

	return sol
}
