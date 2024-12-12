package main

import (
	"fmt"
	utils "github.com/ghjm/advent_utils"
	"github.com/ghjm/advent_utils/board"
	"os"
)

type data struct {
	board *board.StdBoard
}

func getSides(region map[utils.StdPoint]struct{}) int {
	sides := 0
	for p := range region {
		_, left := region[utils.StdPoint{p.X - 1, p.Y}]
		_, right := region[utils.StdPoint{p.X + 1, p.Y}]
		_, up := region[utils.StdPoint{p.X, p.Y - 1}]
		_, down := region[utils.StdPoint{p.X, p.Y + 1}]
		_, upperLeft := region[utils.StdPoint{p.X - 1, p.Y - 1}]
		_, upperRight := region[utils.StdPoint{p.X + 1, p.Y - 1}]
		_, downLeft := region[utils.StdPoint{p.X - 1, p.Y + 1}]
		_, downRight := region[utils.StdPoint{p.X + 1, p.Y + 1}]
		if !left && !up {
			sides++
		}
		if !right && !up {
			sides++
		}
		if !left && !down {
			sides++
		}
		if !right && !down {
			sides++
		}
		if !upperRight && up && right {
			sides++
		}
		if !upperLeft && up && left {
			sides++
		}
		if !downLeft && down && left {
			sides++
		}
		if !downRight && down && right {
			sides++
		}
	}
	return sides
}

func run() error {
	d := data{
		board: board.NewStdBoard(),
	}
	err := d.board.FromFile("input12.txt")
	if err != nil {
		return err
	}
	regions := d.board.FindRegions(false)
	var perimeters []int
	part1 := 0
	for _, reg := range regions {
		perim := 0
		for p := range reg {
			for _, cp := range d.board.Cardinals(p, true) {
				if _, ok := reg[cp]; !ok {
					perim++
				}
			}
		}
		perimeters = append(perimeters, perim)
		part1 += len(reg) * perim
	}
	fmt.Printf("Part 1: %d\n", part1)
	part2 := 0
	for _, reg := range regions {
		part2 += len(reg) * getSides(reg)
	}
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
