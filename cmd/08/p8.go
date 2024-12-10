package main

import (
	"fmt"
	"github.com/ghjm/advent_utils/board"
	"os"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	board *board.RunePlusBoard[int, map[rune]struct{}]
}

func (d *data) getAntinodes(harmonics bool) {
	frequencies := make(map[rune][]utils.Point[int])
	d.board.IterateRunes(func(p utils.Point[int], v rune) bool {
		if v == '.' {
			return true
		}
		_, ok := frequencies[v]
		if !ok {
			frequencies[v] = make([]utils.Point[int], 0)
		}
		frequencies[v] = append(frequencies[v], p)
		return true
	})
	for id, freq := range frequencies {
		for i := 0; i < len(freq); i++ {
			for j := 0; j < len(freq); j++ {
				if i == j {
					continue
				}
				if harmonics {
					an := freq[i]
					for {
						if d.board.Contains(an) {
							if d.board.GetExtra(an) == nil {
								d.board.SetExtra(an, make(map[rune]struct{}))
							}
							d.board.GetExtra(an)[id] = struct{}{}
						} else {
							break
						}
						an = an.Add(freq[i].Delta(freq[j]))
					}
				} else {
					an := freq[i].Add(freq[i].Delta(freq[j]))
					if d.board.Contains(an) {
						if d.board.GetExtra(an) == nil {
							d.board.SetExtra(an, make(map[rune]struct{}))
						}
						d.board.GetExtra(an)[id] = struct{}{}
					}
				}
			}
		}
	}
}

func (d *data) countAntinodeLocs() int {
	result := 0
	d.board.Iterate(func(p utils.Point[int], v board.RunePlusData[map[rune]struct{}]) bool {
		if len(v.Extra) > 0 {
			result++
		}
		return true
	})
	return result
}

func run() error {
	d := data{
		board: board.NewRunePlusBoard[int, map[rune]struct{}](),
	}
	err := d.board.FromFile("input8.txt")
	if err != nil {
		return err
	}
	d.getAntinodes(false)
	fmt.Printf("Part 1: %d\n", d.countAntinodeLocs())
	d.getAntinodes(true)
	fmt.Printf("Part 2: %d\n", d.countAntinodeLocs())
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
