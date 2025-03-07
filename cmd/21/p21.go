package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ghjm/advent_utils"
)

type data struct {
	lines []string
}

var foo string = "foo"

var keyPositions = map[string]utils.StdPoint{
	"0": {1, 3},
	"1": {0, 2},
	"2": {1, 2},
	"3": {2, 2},
	"4": {0, 1},
	"5": {1, 1},
	"6": {2, 1},
	"7": {0, 0},
	"8": {1, 0},
	"9": {2, 0},
	"A": {2, 3},
}

func DoorSequence(start, end string) string {
	order := []string{"<", "v", "^", ">"}
	switch {
	case start == end:
		return "A"
	case (strings.Contains("0A", start) && strings.Contains("147", end)) || (strings.Contains("0A", end) && strings.Contains("147", start)):
		order = []string{"^", ">", "v", "<"}
	}

	startPt := keyPositions[start]
	endPt := keyPositions[end]
	stepsUpDown := endPt.Y - startPt.Y
	stepsLeftRight := endPt.X - startPt.X
	sequence := ""

	for _, i := range order {
		switch {
		case i == "<" && stepsLeftRight < 0:
			sequence += strings.Repeat("<", utils.Abs(stepsLeftRight))
		case i == ">" && stepsLeftRight > 0:
			sequence += strings.Repeat(">", utils.Abs(stepsLeftRight))
		case i == "^" && stepsUpDown < 0:
			sequence += strings.Repeat("^", utils.Abs(stepsUpDown))
		case i == "v" && stepsUpDown > 0:
			sequence += strings.Repeat("v", utils.Abs(stepsUpDown))
		}
	}

	return sequence + "A"
}

var keyPadSequenceMap = map[string]map[string]string{
	"<": {
		">": ">>A",
		"A": ">>^A",
		"^": ">^A",
		"v": ">A",
	},
	">": {
		"<": "<<A",
		"A": "^A",
		"v": "<A",
		"^": "<^A",
	},
	"^": {
		"v": "vA",
		"A": ">A",
		">": "v>A",
		"<": "v<A",
	},
	"v": {
		"^": "^A",
		"A": "^>A",
		"<": "<A",
		">": ">A",
	},
	"A": {
		"<": "v<<A",
		">": "vA",
		"^": "<A",
		"v": "<vA",
	},
}

func KeyPadSequenceCost(start, end string, nRobots int) int {
	var seq string
	var ok bool
	if start == end {
		seq = "A"
		ok = true
	} else {
		seq, ok = keyPadSequenceMap[start][end]
	}

	if ok {
		if nRobots == 1 {
			return len(seq)
		} else {
			return Cost(seq, nRobots-1)
		}
	}

	log.Fatalf("Unknown start %s or end %s", start, end)
	return 0
}

var costCache sync.Map

type costKey struct {
	code    string
	nRobots int
}

func Cost(code string, nRobots int) int {
	cacheKey := costKey{code, nRobots}
	if val, ok := costCache.Load(cacheKey); ok {
		return val.(int)
	}
	result := 0
	code = "A" + code
	for i := 0; i < len(code)-1; i++ {
		// use i:i+1 to get a string instead of just i to get a rune
		result += KeyPadSequenceCost(code[i:i+1], code[i+1:i+2], nRobots)
	}

	costCache.Store(cacheKey, result)
	return result
}

func Complexity(code string, nRobots int) int {
	sequence := ""
	// Append A as start position to code
	code = "A" + code
	for i := 0; i < len(code)-1; i++ {
		sequence += DoorSequence(string(code[i]), string(code[i+1]))
	}

	// recursively find the cost of the shortest sequence
	cost := Cost(sequence, nRobots)
	// log.Println(sequence)
	num, err := strconv.Atoi(code[1 : len(code)-1])
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return cost * num
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input21.txt", func(s string) error {
		d.lines = append(d.lines, s)
		return nil
	})
	if err != nil {
		return err
	}
	sum := 0
	for _, code := range d.lines {
		result := Complexity(code, 2)
		sum += result
	}
	fmt.Printf("Part 1: %d\n", sum)
	sum = 0
	for _, code := range d.lines {
		result := Complexity(code, 25)
		sum += result
	}
	fmt.Printf("Part 2: %d\n", sum)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
