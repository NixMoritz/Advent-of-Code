package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	solvePart1()
	solvePart2()
}

func solvePart1() {
	total := 0
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	rules, toProduce := parseFile(file)
	rulesMap := createRulesMap(rules)

	for _, item := range toProduce {
		working := strings.Split(item, ",")
		if isSafe(working, rulesMap) {
			middleVal, _ := strconv.Atoi(working[len(working)/2])
			total += middleVal
		}
	}

	fmt.Println(total)
}

func solvePart2() {
	total := 0
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	rules, toProduce := parseFile(file)
	rulesMap := createRulesMap(rules)

	for _, item := range toProduce {
		working := strings.Split(item, ",")
		if !isSafe(working, rulesMap) {
			working = reorderValues(working, rulesMap)
			middleVal, _ := strconv.Atoi(working[len(working)/2])
			total += middleVal
		}
	}

	fmt.Println(total)
}

func parseFile(file *os.File) ([]string, []string) {
	var rules, toProduce []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if strings.Contains(line, "|") {
				rules = append(rules, line)
			} else {
				toProduce = append(toProduce, line)
			}
		}
	}

	return rules, toProduce
}

func createRulesMap(rules []string) map[string][]string {
	rulesMap := make(map[string][]string)

	for _, rule := range rules {
		parts := strings.Split(rule, "|")
		if len(parts) == 2 {
			rulesMap[parts[1]] = append(rulesMap[parts[1]], parts[0])
		}
	}

	return rulesMap
}

func isSafe(items []string, rulesMap map[string][]string) bool {
	for i := 0; i < len(items)-1; i++ {
		after := items[i+1:]
		testVals := rulesMap[items[i]]
		for _, val := range after {
			if slices.Contains(testVals, val) {
				return false
			}
		}
	}
	return true
}

func reorderValues(vals []string, rulesMap map[string][]string) []string {
	changed := false
	var reordered []string

	for i, curVal := range vals {
		thingsBefore := rulesMap[curVal]
		thingsAfter := vals[i+1:]

		for _, afterVal := range thingsAfter {
			if slices.Contains(thingsBefore, afterVal) && !slices.Contains(reordered, afterVal) {
				reordered = append(reordered, afterVal)
				changed = true
			}
		}

		if !slices.Contains(reordered, curVal) {
			reordered = append(reordered, curVal)
		}
	}

	if changed {
		return reorderValues(reordered, rulesMap)
	}
	return reordered
}
