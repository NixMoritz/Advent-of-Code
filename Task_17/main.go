package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Computer struct {
	A, B, C uint64
	Program []byte
}

type OutputBuffer []byte

func (buf OutputBuffer) String() string {
	var builder strings.Builder
	for i, b := range buf {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%d", b))
	}
	return builder.String()
}

func (c *Computer) combo(operand byte) uint64 {
	switch operand {
	case 0, 1, 2, 3:
		return uint64(operand)
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	default:
		panic(fmt.Sprintf("Invalid operand: %d", operand))
	}
}

func (c *Computer) Run() OutputBuffer {
	var output OutputBuffer
	for pointer := 0; pointer < len(c.Program)-1; pointer += 2 {
		op, arg := c.Program[pointer], c.Program[pointer+1]
		switch op {
		case 0: // adv
			c.A >>= c.combo(arg)
		case 1: // bxl
			c.B ^= uint64(arg)
		case 2: // bst
			c.B = c.combo(arg) & 7
		case 3: // jnz
			if c.A != 0 {
				pointer = int(arg) - 2
				continue
			}
		case 4: // bxc
			c.B ^= c.C
		case 5: // out
			output = append(output, byte(c.combo(arg)&7))
		case 6: // bdv
			c.B = c.A >> c.combo(arg)
		case 7: // cdv
			c.C = c.A >> c.combo(arg)
		default:
			panic(fmt.Sprintf("Unknown operation: %d", op))
		}
	}
	return output
}

func loadComputerFromFile(filePath string) (*Computer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	computer := &Computer{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Register A:") {
			_, err := fmt.Sscanf(line, "Register A: %d", &computer.A)
			if err != nil {
				return nil, fmt.Errorf("error parsing register A in file '%s': %w", filePath, err)
			}
		} else if strings.HasPrefix(line, "Register B:") {
			_, err := fmt.Sscanf(line, "Register B: %d", &computer.B)
			if err != nil {
				return nil, fmt.Errorf("error parsing register B in file '%s': %w", filePath, err)
			}
		} else if strings.HasPrefix(line, "Register C:") {
			_, err := fmt.Sscanf(line, "Register C: %d", &computer.C)
			if err != nil {
				return nil, fmt.Errorf("error parsing register C in file '%s': %w", filePath, err)
			}
		} else if strings.HasPrefix(line, "Program:") {
			programStr := strings.TrimPrefix(line, "Program: ")
			parts := strings.Split(programStr, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}

				byteVal, err := strconv.ParseUint(part, 10, 8)
				if err != nil {
					return nil, fmt.Errorf("error parsing program byte '%s' in file '%s': %w", part, filePath, err)
				}
				computer.Program = append(computer.Program, byte(byteVal))

			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	return computer, nil
}

func part1(c *Computer) string {
	return c.Run().String()
}

func part2(c *Computer) uint64 {
	var res uint64
	for i := len(c.Program) - 1; i >= 0; i-- {
		res <<= 3
		for {
			c.A, c.B, c.C = uint64(res), 0, 0
			buf := c.Run()
			if !slices.Equal(buf, c.Program[i:]) {
				res++
			} else {
				break
			}
		}
	}
	return res
}

func runPart(part int, filename string) {
	computer, err := loadComputerFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading computer from '%s': %v\n", filename, err)
		return
	}

	if part == 1 {
		result := part1(computer)
		fmt.Printf("Part 1 (%s): %s\n", filename, result)
	} else if part == 2 {
		result := part2(computer)
		fmt.Printf("Part 2 (%s): %d\n", filename, result)
	}
}

func main() {
	runPart(1, "example1.txt")
	runPart(2, "example2.txt")
	runPart(1, "input.txt")
	runPart(2, "input.txt")
}
