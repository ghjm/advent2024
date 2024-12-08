package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	reports [][]int
}

func monotonic(r []int) bool {
	if len(r) < 2 {
		return true
	}
	if r[0] == r[1] {
		return false
	}
	incr := true
	if r[0] > r[1] {
		incr = false
	}
	for i := 1; i < len(r); i++ {
		if (incr && r[i] <= r[i-1]) || (!incr && r[i] >= r[i-1]) {
			return false
		}
	}
	return true
}

func safe(r []int) bool {
	if len(r) < 2 {
		return true
	}
	if !monotonic(r) {
		return false
	}
	for i := 1; i < len(r); i++ {
		diff := utils.Abs(r[i] - r[i-1])
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}

func dampenedSafe(r []int) bool {
	if safe(r) {
		return true
	}
	for i := 0; i < len(r); i++ {
		var newR []int
		for j, v := range r {
			if j != i {
				newR = append(newR, v)
			}
		}
		if safe(newR) {
			return true
		}
	}
	return false
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input2.txt", func(s string) error {
		l := strings.Split(s, " ")
		var report []int
		for _, r := range l {
			report = append(report, utils.MustAtoi(r))
		}
		d.reports = append(d.reports, report)
		return nil
	})
	if err != nil {
		return err
	}
	safeCount := 0
	dampenedSafeCount := 0
	for _, r := range d.reports {
		if safe(r) {
			safeCount++
		}
		if dampenedSafe(r) {
			dampenedSafeCount++
		}
	}
	fmt.Printf("Part 1: %d\n", safeCount)
	fmt.Printf("Part 2: %d\n", dampenedSafeCount)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
