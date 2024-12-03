package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	fmt.Println(solve())
}

func solve() (int, error) {
	inputData := `
		`

	patternMul := `mul\((\d+),(\d+)\)`
	patternDo := `do\(\)`
	patternDont := `don't\(\)`

	reMul := regexp.MustCompile(patternMul)
	reDo := regexp.MustCompile(patternDo)
	reDont := regexp.MustCompile(patternDont)

	matchesMul := reMul.FindAllStringSubmatchIndex(inputData, -1)
	result := 0

	if matchesMul != nil {
		fmt.Println("Found occurrences of mul(digit,digit):")
		for _, match := range matchesMul {

			startMulIndex := match[0]
			num1, err1 := strconv.Atoi(inputData[match[2]:match[3]])
			num2, err2 := strconv.Atoi(inputData[match[4]:match[5]])

			if err1 != nil || err2 != nil {
				return 0, fmt.Errorf("failed to convert string to integer")
			}

			substringBeforeMul := inputData[:startMulIndex]

			allDoMatches := reDo.FindAllStringIndex(substringBeforeMul, -1)
			allDontMatches := reDont.FindAllStringIndex(substringBeforeMul, -1)

			lastDoMatch := getLastMatch(allDoMatches)
			lastDontMatch := getLastMatch(allDontMatches)

			fmt.Println("Last do() match:", lastDoMatch)
			fmt.Println("Last don't() match:", lastDontMatch)

			if lastDontMatch != nil && (lastDoMatch == nil || lastDontMatch[0] < lastDoMatch[0]) {
				fmt.Printf("Skipping mul(%d, %d) due to previous 'don't()'\n", num1, num2)
				continue
			}

			if lastDoMatch != nil {
				fmt.Printf("Calculating mul(%d, %d) due to previous 'do()'\n", num1, num2)
				result += mul(num1, num2)
			}
		}
	} else {
		fmt.Println("No mul matches found.")
	}

	return result, nil
}

func getLastMatch(matches [][]int) []int {
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return nil
}

func mul(no1 int, no2 int) int {
	return no1 * no2
}
