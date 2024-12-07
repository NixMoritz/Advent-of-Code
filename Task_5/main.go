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
	p1()
	p2()
}

func p1() {
	var total int = 0
	file, _ := os.Open("input.txt")
	defer file.Close()

	var rules []string
	var toProduce []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.Contains(line, "|") {
			rules = append(rules, line)
		} else if len(line) > 0 {
			toProduce = append(toProduce, line)
		}
	}

	var rulesMap map[string][]string = make(map[string][]string)

	for i := 0; i < len(rules); i++ {
		working := strings.Split(rules[i], "|")
		rulesMap[working[1]] = append(rulesMap[working[1]], working[0])
	}

	for i := 0; i < len(toProduce); i++ {
		working := strings.Split(toProduce[i], ",")
		var safe bool = true
		for j := 0; j < len(working)-1; j++ {
			after := working[j+1:]
			testVals := rulesMap[working[j]]
			for k := 0; k < len(after); k++ {
				if slices.Contains(testVals, after[k]) {
					safe = false
				}
			}
		}
		if safe {
			middleVal, _ := strconv.Atoi(working[len(working)/2])
			total += middleVal
		}
	}
	fmt.Println(total)
}

func p2() {
	var total int = 0
	file, _ := os.Open("input.txt")
	defer file.Close()

	var rules []string
	var toProduce []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.Contains(line, "|") {
			rules = append(rules, line)
		} else if len(line) > 0 {
			toProduce = append(toProduce, line)
		}
	}

	var rulesMap map[string][]string = make(map[string][]string)

	for i := 0; i < len(rules); i++ {
		working := strings.Split(rules[i], "|")
		rulesMap[working[1]] = append(rulesMap[working[1]], working[0])
	}

	for i := 0; i < len(toProduce); i++ {
		working := strings.Split(toProduce[i], ",")
		var safe bool = true
		for j := 0; j < len(working)-1; j++ {
			after := working[j+1:]
			testVals := rulesMap[working[j]]
			for k := 0; k < len(after); k++ {
				if slices.Contains(testVals, after[k]) {
					safe = false
				}
			}
		}
		if !safe {
			working = reorderVals(working, rulesMap)
			middleVal, _ := strconv.Atoi(working[len(working)/2])
			total += middleVal
		}
	}
	fmt.Println(total)
}

func reorderVals(vals []string, rules map[string][]string) []string {
	var changed bool = false
	var temp []string
	for i := 0; i < len(vals); i++ {
		curVal := vals[i]
		thingsShouldBeBefore := rules[curVal]
		thingsAfter := vals[i+1:]

		for j := 0; j < len(thingsAfter); j++ {
			if slices.Contains(thingsShouldBeBefore, thingsAfter[j]) && !slices.Contains(temp, thingsAfter[j]) {
				temp = append(temp, thingsAfter[j])
				changed = true
			}
		}
		if !slices.Contains(temp, curVal) {
			temp = append(temp, curVal)
		}
	}
	if !changed {
		return temp
	} else {
		return reorderVals(temp, rules)
	}
}
