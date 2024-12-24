package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Gate struct {
	Op     string
	Input1 string
	Input2 string
	Output string
}

type Puzzle struct {
	WireValues map[string]int
	Operations []Operation
	ZBits      int
}

type Operation struct {
	Res, A, Op, B string
}

func main() {
	fmt.Println("Part 1:")
	exampleResult, err := runSimulation("example.txt", nil)
	if err != nil {
		fmt.Printf("Failed to simulate example: %v\n", err)
		return
	}
	fmt.Printf("Example Result: %d (should be 2024)\n", exampleResult)

	result, err := runSimulation("input.txt", nil)
	if err != nil {
		fmt.Printf("Failed to simulate input: %v\n", err)
		return
	}
	fmt.Printf("Part 1 Result: %d\n\n", result)

	fmt.Println("Part 2:")
	puzzle, err := readInputFile("input.txt")
	if err != nil {
		fmt.Printf("Failed to read puzzle input: %v\n", err)
		return
	}
	solution := solution(puzzle)
	slog.Info(fmt.Sprint("Solution is ", solution))
}

func runSimulation(filePath string, swaps map[string]string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return simulateGates(string(content), swaps), nil
}

func simulateGates(input string, swaps map[string]string) int {
	wires := make(map[string]int)
	gates := parseInput(input, wires, swaps)
	evaluateCircuit(wires, gates)
	return calculateZValue(wires)
}

func parseInput(input string, wires map[string]int, swaps map[string]string) []Gate {
	var gates []Gate
	scanner := bufio.NewScanner(strings.NewReader(input))
	parsingInitialValues := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			parsingInitialValues = false
			continue
		}

		if parsingInitialValues {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				value := 0
				if parts[1] == "1" {
					value = 1
				}
				wires[parts[0]] = value
			}
		} else {
			parts := strings.Fields(line)
			if len(parts) == 5 {
				output := parts[4]
				if swaps != nil {
					if swapped, ok := swaps[output]; ok {
						output = swapped
					}
				}
				gates = append(gates, Gate{
					Input1: parts[0],
					Op:     parts[1],
					Input2: parts[2],
					Output: output,
				})
			}
		}
	}
	return gates
}

func evaluateCircuit(wires map[string]int, gates []Gate) {
	changed := true
	for changed {
		changed = false
		for _, gate := range gates {
			if _, exists := wires[gate.Output]; exists {
				continue
			}

			val1, ok1 := wires[gate.Input1]
			val2, ok2 := wires[gate.Input2]
			if !ok1 || !ok2 {
				continue
			}

			result := evaluateGate(gate, val1, val2)
			wires[gate.Output] = result
			changed = true
		}
	}
}

func evaluateGate(gate Gate, val1, val2 int) int {
	switch gate.Op {
	case "AND":
		return val1 & val2
	case "OR":
		return val1 | val2
	case "XOR":
		return val1 ^ val2
	default:
		return 0
	}
}

func calculateZValue(wires map[string]int) int {
	var zWires []string
	for wire := range wires {
		if strings.HasPrefix(wire, "z") {
			zWires = append(zWires, wire)
		}
	}
	sort.Strings(zWires)

	result := 0
	for i := len(zWires) - 1; i >= 0; i-- {
		result = (result << 1) | wires[zWires[i]]
	}
	return result
}

func calculateWireValue(wires map[string]int, prefix string) int {
	var filteredWires []string
	for wire := range wires {
		if strings.HasPrefix(wire, prefix) {
			filteredWires = append(filteredWires, wire)
		}
	}
	sort.Strings(filteredWires)

	result := 0
	for i := len(filteredWires) - 1; i >= 0; i-- {
		result = (result << 1) | wires[filteredWires[i]]
	}
	return result
}

func solution(puzzle *Puzzle) string {
	andGatesOfInputs := make(map[int]string)
	xorGatesOfInputs := make(map[int]string)
	producedBy := make(map[string]Operation)
	swappedGates := make(map[string]int)

	for _, op := range puzzle.Operations {
		producedBy[op.Res] = op

		if (strings.HasPrefix(op.A, "x") && strings.HasPrefix(op.B, "y")) || (strings.HasPrefix(op.B, "x") && strings.HasPrefix(op.A, "y")) {
			bitsA, _ := strconv.Atoi(op.A[1:])
			bitsB, _ := strconv.Atoi(op.B[1:])
			if bitsA != bitsB {
				panic(fmt.Sprintf("Using inputs with different bit-index %+v\n", op))
			}

			if op.Op == "OR" {
				panic(fmt.Sprintf("Using inputs with an OR %+v\n", op))
			}

			if op.Op == "AND" {
				andGatesOfInputs[bitsA] = op.Res
			} else {
				xorGatesOfInputs[bitsA] = op.Res
			}
		}
	}

	for bit := range puzzle.ZBits - 1 {
		if _, ok := andGatesOfInputs[bit]; !ok {
			panic(fmt.Sprintf("Missing andGatesOfInput for %d\n", bit))
		}
		if _, ok := xorGatesOfInputs[bit]; !ok {
			panic(fmt.Sprintf("Missing xorGatesOfInputs for %d\n", bit))
		}
	}

	{
		op := producedBy["z00"]
		if op.Op != "XOR" || op.A != "x00" || op.B != "y00" {
			swappedGates["z00"] = 1
		}
	}

	{
		name := fmt.Sprintf("z%02d", puzzle.ZBits-1)
		op, ok := producedBy[name]

		if !ok {
			panic("missing last z")
		}

		if op.Op != "OR" {
			swappedGates[op.Res] = 1
		} else {
			left := producedBy[op.A]
			right := producedBy[op.B]

			if left.Op != "AND" || right.Op != "AND" {
				panic("good luck")
			}
		}
	}

	for bit := 1; bit < puzzle.ZBits-1; bit++ {
		name := fmt.Sprintf("z%02d", bit)

		op, ok := producedBy[name]
		if !ok {
			panic(fmt.Sprintf("missing %s", name))
		}

		if op.Op != "XOR" {
			swappedGates[op.Res] = 1
			continue
		}

		myXorIn := xorGatesOfInputs[bit]
		left := producedBy[op.A]
		right := producedBy[op.B]

		if right.Res == myXorIn {
			left, right = right, left
		}

		leftOk := true
		if left.Res != myXorIn && left.Res != "" {
			swappedGates[myXorIn] = 1
			leftOk = false
		}

		rightOk := true
		if right.Op != "OR" {
			if !(bit == 1 && right.Res == andGatesOfInputs[0]) {
				swappedGates[right.Res] = 1
				rightOk = false
			}
		}

		if !leftOk && rightOk {
			swappedGates[left.Res] = 1
		}

		if !rightOk || right.Op != "OR" {
			continue
		}

		left = producedBy[right.A]
		right = producedBy[right.B]

		if left.Op != "AND" {
			swappedGates[left.Res] = 1
		}

		if right.Op != "AND" {
			swappedGates[right.Res] = 1
		}
	}

	swappedWireNames := make([]string, 0, len(swappedGates))
	for name := range swappedGates {
		swappedWireNames = append(swappedWireNames, name)
	}
	slices.Sort(swappedWireNames)
	return strings.Join(swappedWireNames, ",")
}

func readScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	wireValues := make(map[string]int)
	operations := make([]Operation, 0)
	zbits := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if strings.Contains(line, ":") {
			inputWireValue := strings.Split(line, ": ")
			value, err := strconv.Atoi(inputWireValue[1])
			if err != nil {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, err
			}
			wireName := inputWireValue[0]
			wireValues[wireName] = value

			if bitsStr := strings.TrimPrefix(wireName, "z"); bitsStr != wireName {
				bits, err := strconv.Atoi(bitsStr)
				if err != nil {
					return &Puzzle{
						Operations: operations, ZBits: zbits, WireValues: wireValues,
					}, fmt.Errorf("invalid Z wire in line %s due to %w", line, err)
				}
				zbits = max(zbits, bits+1)
			}
		} else {
			var a, b, op, res string
			n, err := fmt.Sscanf(line, "%s %s %s -> %s", &a, &op, &b, &res)
			if err != nil {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, err
			}
			if n != 4 {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, fmt.Errorf("parsed %d/4 elements from line %s", n, line)
			}

			if a > b {
				a, b = b, a
			}
			operations = append(operations, Operation{A: a, B: b, Op: op, Res: res})

			if bitsStr := strings.TrimPrefix(res, "z"); bitsStr != res {
				bits, err := strconv.Atoi(bitsStr)
				if err != nil {
					return &Puzzle{
						Operations: operations, ZBits: zbits, WireValues: wireValues,
					}, fmt.Errorf("invalid Z wire in line %s due to %w", line, err)
				}
				zbits = max(zbits, bits+1)
			}
		}
	}
	return &Puzzle{
		Operations: operations, ZBits: zbits, WireValues: wireValues,
	}, scanner.Err()
}

func readInputFile(path string) (*Puzzle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	return readScanner(scanner)
}
