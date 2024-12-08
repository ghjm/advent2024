package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type equation struct {
	value     uint64
	numbers   []uint64
	canConcat bool
}

type data struct {
	eqs []equation
}

func (eq *equation) tryEval(curValue uint64, remaining []uint64) uint64 {
	if len(remaining) == 0 || curValue > eq.value {
		return curValue
	}
	if eq.tryEval(curValue+remaining[0], remaining[1:]) == eq.value {
		return eq.value
	}
	if eq.tryEval(curValue*remaining[0], remaining[1:]) == eq.value {
		return eq.value
	}
	if eq.canConcat && eq.tryEval(utils.MustAtoiU64(strconv.FormatUint(curValue, 10)+strconv.FormatUint(remaining[0], 10)), remaining[1:]) == eq.value {
		return eq.value
	}
	return eq.value + 1
}

func (eq *equation) CanEval() bool {
	if len(eq.numbers) == 0 {
		return false
	}
	if len(eq.numbers) == 1 {
		return eq.numbers[0] == eq.value
	}
	return eq.tryEval(eq.numbers[0], eq.numbers[1:]) == eq.value
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input7.txt", func(s string) error {
		var eq equation
		v1 := strings.Split(s, ": ")
		if len(v1) != 2 {
			return fmt.Errorf("invalid input")
		}
		eq.value = utils.MustAtoiU64(v1[0])
		for _, v2 := range strings.Split(v1[1], " ") {
			eq.numbers = append(eq.numbers, utils.MustAtoiU64(v2))
		}
		d.eqs = append(d.eqs, eq)
		return nil
	})
	if err != nil {
		return err
	}
	var sum1, sum2 uint64
	for _, eq := range d.eqs {
		if eq.CanEval() {
			sum1 += eq.value
			sum2 += eq.value
		} else {
			eq.canConcat = true
			if eq.CanEval() {
				sum2 += eq.value
			}
		}
	}
	fmt.Printf("Part 1: %d\n", sum1)
	fmt.Printf("Part 2: %d\n", sum2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
