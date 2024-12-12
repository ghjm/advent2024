package main

import (
	"fmt"
	"github.com/ghjm/advent_utils/board"
	"os"

	"github.com/ghjm/advent_utils"
)

type data struct {
	board *board.StdBoard
}

func (d *data) score(start utils.StdPoint) int {
	open := []utils.StdPoint{start}
	visited := make(map[utils.StdPoint]struct{})
	summits := make(map[utils.StdPoint]struct{})
	for len(open) > 0 {
		p := open[0]
		open = open[1:]
		if _, ok := visited[p]; ok {
			continue
		}
		visited[p] = struct{}{}
		pv := d.board.Get(p)
		for _, c := range d.board.Cardinals(p, false) {
			if _, ok := visited[c]; ok {
				continue
			}
			cv := d.board.Get(c)
			if cv == pv+1 {
				if cv == '9' {
					summits[c] = struct{}{}
				} else {
					open = append(open, c)
				}
			}
		}
	}
	return len(summits)
}

func (d *data) rating(p utils.StdPoint) int {
	pv := d.board.Get(p)
	if pv == '9' {
		return 1
	}
	sum := 0
	for _, c := range d.board.Cardinals(p, false) {
		cv := d.board.Get(c)
		if cv == pv+1 {
			sum += d.rating(c)
		}
	}
	return sum
}

func run() error {
	d := data{
		board: board.NewStdBoard(board.WithStorage(&board.FlatBoard{})),
	}
	err := d.board.FromFile("input10.txt")
	if err != nil {
		return err
	}
	part1 := 0
	part2 := 0
	d.board.Iterate(func(p utils.Point[int], v rune) bool {
		if v == '0' {
			part1 += d.score(p)
			part2 += d.rating(p)
		}
		return true
	})
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
