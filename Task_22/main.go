package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// processSecret performs one complete iteration of the secret number transformation
func processSecret(secret int) int {
	result := secret * 64
	secret ^= result   // mix
	secret %= 16777216 // prune

	result = secret / 32
	secret ^= result   // mix
	secret %= 16777216 // prune

	result = secret * 2048
	secret ^= result   // mix
	secret %= 16777216 // prune

	return secret
}

// generateNthSecret generates the nth new secret number for a given initial secret
func generateNthSecret(initialSecret, n int) int {
	secret := initialSecret
	for i := 0; i < n; i++ {
		secret = processSecret(secret)
	}
	return secret
}

// readInput reads the input from the specified file and converts it to a slice of integers
func readInput(filename string) ([]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	secrets := make([]int, 0, len(lines))
	for _, line := range lines {
		secret, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			continue
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// Part 1: Calculate the sum of the 2000th secret for all buyers
func part1(secrets []int) int {
	sum := 0
	for _, secret := range secrets {
		secret2000 := generateNthSecret(secret, 2000)
		sum += secret2000
	}
	return sum
}

// generatePricesAndChanges generates all prices and changes for a secret
func generatePricesAndChanges(initialSecret int, n int) ([]int, map[string]int) {
	prices := make([]int, n+1)
	secret := initialSecret
	prices[0] = secret % 10

	for i := 1; i <= n; i++ {
		secret = processSecret(secret)
		prices[i] = secret % 10
	}

	sequences := make(map[string]int)
	for i := 3; i < n; i++ {
		changes := []int{
			prices[i-2] - prices[i-3],
			prices[i-1] - prices[i-2],
			prices[i] - prices[i-1],
			prices[i+1] - prices[i],
		}
		key := fmt.Sprintf("%d,%d,%d,%d", changes[0], changes[1], changes[2], changes[3])
		if _, exists := sequences[key]; !exists {
			sequences[key] = prices[i+1]
		}
	}

	return prices, sequences
}

// Part 2: Find the best sequence among all secrets in the list
func part2(secrets []int) (string, int) {
	totalBananas := make(map[string]int)

	for _, secret := range secrets {
		_, sequences := generatePricesAndChanges(secret, 2000)
		for seq, price := range sequences {
			totalBananas[seq] += price
		}
	}

	bestSequence := ""
	maxBananas := 0
	for seq, bananas := range totalBananas {
		if bananas > maxBananas {
			maxBananas = bananas
			bestSequence = seq
		}
	}

	return bestSequence, maxBananas
}

func main() {
	// Read input
	secrets, err := readInput("input.txt")
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Part 1
	fmt.Println("Running Part 1...")
	sum := part1(secrets)
	fmt.Printf("Part 1: The sum of the 2000th secret number for all buyers is %d\n", sum)

	// Part 2
	fmt.Println("\nRunning Part 2...")
	bestSequence, totalBananas := part2(secrets)
	fmt.Printf("Part 2: Best sequence is %s with a total of %d bananas\n", bestSequence, totalBananas)
}
