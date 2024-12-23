package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// Read input from file
	b, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Split input into lines
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")

	// Get solutions
	solution1 := findTriples(lines)
	solution2 := findLargestClique(lines)

	fmt.Printf("Part 1 Solution: %d\n", solution1)
	fmt.Printf("Part 2 Solution: %s\n", strings.Join(solution2, ","))
}

// findTriples returns all sets of three interconnected computers where at least
// one computer name starts with 't'
func findTriples(connections []string) int {
	// Build adjacency map
	adj := make(map[string]map[string]bool)

	// Parse connections and build adjacency map
	for _, conn := range connections {
		parts := strings.Split(strings.TrimSpace(conn), "-")
		a, b := parts[0], parts[1]

		// Initialize maps if they don't exist
		if adj[a] == nil {
			adj[a] = make(map[string]bool)
		}
		if adj[b] == nil {
			adj[b] = make(map[string]bool)
		}

		// Add bidirectional connections
		adj[a][b] = true
		adj[b][a] = true
	}

	// Get all unique computer names
	computers := make([]string, 0)
	for comp := range adj {
		computers = append(computers, comp)
	}

	// Find all triples
	count := 0

	for i := 0; i < len(computers); i++ {
		for j := i + 1; j < len(computers); j++ {
			// Check if i and j are connected
			if !adj[computers[i]][computers[j]] {
				continue
			}

			for k := j + 1; k < len(computers); k++ {
				// Check if all three computers are interconnected
				if adj[computers[i]][computers[k]] && adj[computers[j]][computers[k]] {
					// Check if at least one computer starts with 't'
					if strings.HasPrefix(computers[i], "t") ||
						strings.HasPrefix(computers[j], "t") ||
						strings.HasPrefix(computers[k], "t") {
						count++
					}
				}
			}
		}
	}

	return count
}

func findLargestClique(connections []string) []string {
	// Build adjacency map
	adj := make(map[string]map[string]bool)

	// Parse connections and build adjacency map
	for _, conn := range connections {
		parts := strings.Split(strings.TrimSpace(conn), "-")
		a, b := parts[0], parts[1]

		// Initialize maps if they don't exist
		if adj[a] == nil {
			adj[a] = make(map[string]bool)
		}
		if adj[b] == nil {
			adj[b] = make(map[string]bool)
		}

		// Add bidirectional connections
		adj[a][b] = true
		adj[b][a] = true
	}

	// Convert adjacency map keys to a sorted list of computers
	computers := make([]string, 0)
	for comp := range adj {
		computers = append(computers, comp)
	}
	sort.Strings(computers)

	// Bron-Kerbosch algorithm implementation
	var maxClique []string
	var bronKerbosch func(r, p, x []string)

	bronKerbosch = func(r, p, x []string) {
		if len(p) == 0 && len(x) == 0 {
			// If R is a maximal clique and larger than maxClique, update maxClique
			if len(r) > len(maxClique) {
				maxClique = make([]string, len(r))
				copy(maxClique, r)
			}
			return
		}

		// Iterate over a copy of P to avoid modifying it during iteration
		for _, v := range append([]string{}, p...) {
			// Create new sets R, P, and X for recursion
			newR := append(r, v)
			newP := filter(p, func(u string) bool { return adj[v][u] })
			newX := filter(x, func(u string) bool { return adj[v][u] })

			bronKerbosch(newR, newP, newX)

			// Move v from P to X
			p = removeString(p, v)
			x = append(x, v)
		}
	}

	// Start Bron-Kerbosch with R = {}, P = computers, X = {}
	bronKerbosch([]string{}, computers, []string{})

	sort.Strings(maxClique)
	return maxClique
}

// Utility function to filter a slice based on a condition
func filter(slice []string, cond func(string) bool) []string {
	result := []string{}
	for _, v := range slice {
		if cond(v) {
			result = append(result, v)
		}
	}
	return result
}

// removeString removes a string from a slice
func removeString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
