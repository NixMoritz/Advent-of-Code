package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BlockType uint8

const (
	File BlockType = iota
	Free
)

type Block struct {
	blockType BlockType
	num       int
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	input := strings.TrimSpace(string(data))

	// Solve Part 1
	blocks, files := initializeBlocks(input)
	part1Result := solvePart1(blocks)
	fmt.Println("Part 1 Result:", part1Result)

	// Solve Part 2
	blocks, files = initializeBlocks(input)
	part2Result := solvePart2(blocks, files)
	fmt.Println("Part 2 Result:", part2Result)
}

func initializeBlocks(input string) ([]Block, int) {
	files := 0
	blocks := make([]Block, 0)

	for i, b := range input {
		size, _ := strconv.Atoi(string(b))
		var block BlockType
		if i%2 == 0 {
			block = File
			files++
		} else {
			block = Free
		}
		for j := 0; j < size; j++ {
			blocks = append(blocks, Block{block, files - 1})
		}
	}

	return blocks, files
}

func solvePart1(blocks []Block) int {
	start := 0
	for i := len(blocks) - 1; i > (len(blocks) / 2); i-- {
		if blocks[i].blockType == File {
			for j := start; j < i; j++ {
				if blocks[j].blockType == Free {
					start = j
					blocks[i], blocks[j] = blocks[j], blocks[i]
					break
				}
			}
		}
	}
	return calculateChecksum(blocks)
}

func solvePart2(blocks []Block, files int) int {
	search := files - 1
	fileBlocks := make([]int, 0)
	freeBlocks := make([]int, 0)

	for i := len(blocks) - 1; i > 0; i-- {
		if blocks[i].blockType == File && blocks[i].num == search {
			fileBlocks = append(fileBlocks, i)
		} else {
			if len(fileBlocks) > 0 {
				for j := 0; j <= i; j++ {
					if blocks[j].blockType == Free {
						freeBlocks = append(freeBlocks, j)
						if len(freeBlocks) == len(fileBlocks) {
							for k := 0; k < len(fileBlocks); k++ {
								blocks[fileBlocks[k]], blocks[freeBlocks[k]] = blocks[freeBlocks[k]], blocks[fileBlocks[k]]
							}
							break
						}
					} else {
						freeBlocks = nil
					}
				}
				freeBlocks = nil
				fileBlocks = nil
				search--
				i++
			}
		}
	}
	return calculateChecksum(blocks)
}

func calculateChecksum(blocks []Block) int {
	checksum := 0
	for i, block := range blocks {
		if block.blockType == File {
			checksum += i * block.num
		}
	}
	return checksum
}
