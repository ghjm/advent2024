package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghjm/advent_utils"
)

type data struct {
	towels          []string
	designs         []string
	makeableDesigns map[string]int
}

func (d *data) waysToMake(design string) int {
	if v, ok := d.makeableDesigns[design]; ok {
		return v
	}
	count := 0
	for _, towel := range d.towels {
		if strings.HasPrefix(design, towel) {
			if len(towel) == len(design) {
				count += 1
			} else {
				count += d.waysToMake(design[len(towel):])
			}
		}
	}
	d.makeableDesigns[design] = count
	return count
}

func run() error {
	d := data{
		makeableDesigns: make(map[string]int),
	}
	err := utils.OpenAndReadLines("input19.txt", func(s string) error {
		if s == "" {
			return nil
		}
		if strings.Contains(s, ",") {
			d.towels = strings.Split(s, ", ")
			return nil
		}
		d.designs = append(d.designs, s)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	part2 := 0
	for _, design := range d.designs {
		wtm := d.waysToMake(design)
		if wtm > 0 {
			part1++
		}
		part2 += wtm
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
