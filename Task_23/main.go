package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	adj := buildAdjacencyMap(lines)

	solution1 := findTriples(adj)
	solution2 := findLargestClique(adj)

	fmt.Printf("Part 1 Solution: %d\n", solution1)
	fmt.Printf("Part 2 Solution: %s\n", strings.Join(solution2, ","))
}

func buildAdjacencyMap(connections []string) map[string]map[string]bool {
	adj := make(map[string]map[string]bool)
	for _, conn := range connections {
		parts := strings.Split(strings.TrimSpace(conn), "-")
		a, b := parts[0], parts[1]

		if adj[a] == nil {
			adj[a] = make(map[string]bool)
		}
		if adj[b] == nil {
			adj[b] = make(map[string]bool)
		}

		adj[a][b] = true
		adj[b][a] = true
	}
	return adj
}

func findTriples(adj map[string]map[string]bool) int {
	computers := make([]string, 0, len(adj))
	for comp := range adj {
		computers = append(computers, comp)
	}

	count := 0
	for i := 0; i < len(computers); i++ {
		for j := i + 1; j < len(computers); j++ {
			if !adj[computers[i]][computers[j]] {
				continue
			}
			for k := j + 1; k < len(computers); k++ {
				if adj[computers[i]][computers[k]] && adj[computers[j]][computers[k]] {
					if strings.HasPrefix(computers[i], "t") || strings.HasPrefix(computers[j], "t") || strings.HasPrefix(computers[k], "t") {
						count++
					}
				}
			}
		}
	}

	return count
}

func findLargestClique(adj map[string]map[string]bool) []string {
	computers := make([]string, 0, len(adj))
	for comp := range adj {
		computers = append(computers, comp)
	}
	sort.Strings(computers)

	var maxClique []string
	var bronKerbosch func(r, p, x []string)

	bronKerbosch = func(r, p, x []string) {
		if len(p) == 0 && len(x) == 0 {
			if len(r) > len(maxClique) {
				maxClique = append([]string(nil), r...)
			}
			return
		}

		for _, v := range append([]string{}, p...) {
			newR := append(r, v)
			newP := filter(p, func(u string) bool { return adj[v][u] })
			newX := filter(x, func(u string) bool { return adj[v][u] })

			bronKerbosch(newR, newP, newX)

			p = removeString(p, v)
			x = append(x, v)
		}
	}

	bronKerbosch([]string{}, computers, []string{})
	sort.Strings(maxClique)
	return maxClique
}

func filter(slice []string, cond func(string) bool) []string {
	result := slice[:0]
	for _, v := range slice {
		if cond(v) {
			result = append(result, v)
		}
	}
	return result
}

func removeString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
