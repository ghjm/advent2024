package main

import (
	"fmt"
	"github.com/ghjm/advent2024/pkg/utils"
	"os"
)

type data struct {
	lines []string
}

type point struct {
	x, y int
}

func (d *data) getAtPos(x, y int) string {
	if x < 0 || y < 0 || y >= len(d.lines) || x >= len(d.lines[y]) {
		return ""
	}
	return d.lines[y][x : x+1]
}

func (d *data) checkXmas(x, y int) int {
	str := "XMAS"
	result := 0
	if d.getAtPos(x, y) != str[0:1] {
		return 0
	}
	for _, dir := range []point{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	} {
		pos := point{x, y}
		count := 0
		for i := 1; i < len(str); i++ {
			pos = point{pos.x + dir.x, pos.y + dir.y}
			count++
			if d.getAtPos(pos.x, pos.y) != str[count:count+1] {
				break
			}
			if count+1 == len(str) {
				result++
			}
		}
	}
	return result
}

func (d *data) checkCross(x, y int) bool {
	if d.getAtPos(x, y) != "A" {
		return false
	}
	l1 := d.getAtPos(x-1, y-1)
	l2 := d.getAtPos(x+1, y-1)
	r1 := d.getAtPos(x-1, y+1)
	r2 := d.getAtPos(x+1, y+1)
	mC := 0
	sC := 0
	for _, s := range []string{l1, l2, r1, r2} {
		if s == "M" {
			mC++
		} else if s == "S" {
			sC++
		} else {
			return false
		}
	}
	if mC != 2 || sC != 2 {
		return false
	}
	for l, r := range map[string]string{l1: r2, l2: r1} {
		if l == r {
			return false
		}
	}
	return true
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input4.txt", func(s string) error {
		d.lines = append(d.lines, s)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	part2 := 0
	for y := 0; y < len(d.lines); y++ {
		for x := 0; x < len(d.lines[y]); x++ {
			part1 += d.checkXmas(x, y)
			if d.checkCross(x, y) {
				part2++
			}
		}
	}
	fmt.Printf("Part one: %d\n", part1)
	fmt.Printf("Part two: %d\n", part2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
