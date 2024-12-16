package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ghjm/advent_utils"
)

type puzzle struct {
	xa, ya, xb, yb, xp, yp int64
}

type data struct {
	puzzles []puzzle
}

func solve(p puzzle) (int64, int64) {
	bn := ((p.ya * p.xp) - (p.xa * p.yp))
	bd := ((p.ya * p.xb) - (p.xa * p.yb))
	if bn%bd != 0 {
		return -1, -1
	}
	b := bn / bd
	a := (p.xp - (p.xb * b)) / p.xa
	return a, b
}

func part2ize(p puzzle) puzzle {
	return puzzle{
		xa: p.xa,
		ya: p.ya,
		xb: p.xb,
		yb: p.yb,
		xp: p.xp + 10000000000000,
		yp: p.yp + 10000000000000,
	}
}

func run() error {
	d := data{}
	s, err := utils.OpenAndReadAll("input13.txt")
	if err != nil {
		return err
	}
	r := regexp.MustCompile(`^Button A: X\+(\d+), Y\+(\d+) Button B: X\+(\d+), Y\+(\d+) Prize: X=(\d+), Y=(\d+)$`)
	for _, p := range strings.Split(string(s), "\n\n") {
		p = strings.Replace(p, "\n", " ", -1)
		m := r.FindStringSubmatch(p)
		if m == nil || len(m) != 7 {
			return fmt.Errorf("invalid puzzle")
		}
		d.puzzles = append(d.puzzles, puzzle{
			xa: utils.MustAtoi64(m[1]),
			ya: utils.MustAtoi64(m[2]),
			xb: utils.MustAtoi64(m[3]),
			yb: utils.MustAtoi64(m[4]),
			xp: utils.MustAtoi64(m[5]),
			yp: utils.MustAtoi64(m[6]),
		})
	}
	var part1 int64
	var part2 int64
	for _, p := range d.puzzles {
		a, b := solve(p)
		if a >= 0 {
			part1 += 3*a + b
		}
		a, b = solve(part2ize(p))
		if a >= 0 {
			part2 += 3*a + b
		}
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
