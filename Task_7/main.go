package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func PartOne(input string) int64 {
	return Sum(Filter(input, Check1))
}

func PartTwo(input string) int64 {
	return Sum(Filter(input, Check2))
}

func Filter(input string, check func(target, acc int64, nums []int64) bool) []int64 {
	var results []int64
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		matches := regexp.MustCompile(`\d+`).FindAllString(line, -1)
		if len(matches) == 0 {
			continue
		}
		var nums []int64
		for _, m := range matches {
			n, _ := strconv.ParseInt(m, 10, 64)
			nums = append(nums, n)
		}
		target := nums[0]
		nums = nums[1:]
		if check(target, 0, nums) {
			results = append(results, target)
		}
	}
	return results
}

func Check1(target, acc int64, nums []int64) bool {
	if len(nums) == 0 {
		return target == acc
	}
	return Check1(target, acc*nums[0], nums[1:]) ||
		Check1(target, acc+nums[0], nums[1:])
}

func Check2(target, acc int64, nums []int64) bool {
	if acc > target {
		return false
	}
	if len(nums) == 0 {
		return target == acc
	}
	concat, _ := strconv.ParseInt(strconv.FormatInt(acc, 10)+strconv.FormatInt(nums[0], 10), 10, 64)
	return Check2(target, concat, nums[1:]) ||
		Check2(target, acc*nums[0], nums[1:]) ||
		Check2(target, acc+nums[0], nums[1:])
}

func Sum(nums []int64) int64 {
	sum := int64(0)
	for _, n := range nums {
		sum += n
	}
	return sum
}

func ReadInputFromFile(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	// Example usage
	input := "15 1 2 3\n20 4 5 6"
	println(PartOne(input))
	println(PartTwo(input))

	// Read input from file
	fileInput, err := ReadInputFromFile("input.txt")
	if err != nil {
		panic(err)
	}
	println(PartOne(fileInput))
	println(PartTwo(fileInput))
}
