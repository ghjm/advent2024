package main

import (
	"fmt"
	"os"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	board          *utils.StandardBoard
	initialPoint   utils.Point[int]
	initialPointer rune
}

var pointers = map[rune]utils.Point[int]{
	'^': {0, -1},
	'v': {0, 1},
	'<': {-1, 0},
	'>': {1, 0},
}

var turnRight = map[rune]rune{
	'^': '>',
	'>': 'v',
	'v': '<',
	'<': '^',
}

type state struct {
	pos     utils.Point[int]
	pointer rune
}

func (d *data) runBoard() (map[utils.Point[int]]struct{}, bool) {
	visited := make(map[state]struct{})
	curPos := d.initialPoint
	curPointer := d.initialPointer
	hasLoop := false
	for {
		visited[state{curPos, curPointer}] = struct{}{}
		nextPos := curPos.Add(pointers[curPointer])
		if !d.board.Contains(nextPos) {
			break
		}
		_, ok := visited[state{nextPos, curPointer}]
		if ok {
			hasLoop = true
			break
		}
		nv := d.board.Get(nextPos)
		if nv != '.' {
			curPointer = turnRight[curPointer]
		} else {
			curPos = nextPos
		}
	}
	vPos := make(map[utils.Point[int]]struct{})
	for k := range visited {
		vPos[k.pos] = struct{}{}
	}
	return vPos, hasLoop
}

func run() error {
	d := data{
		board: utils.NewStandardBoard(),
	}
	d.board.MustFromFile("input6.txt")
	d.board.Transform(func(p utils.Point[int], v rune) rune {
		_, ok := pointers[v]
		if ok {
			if d.initialPointer != 0 {
				panic(fmt.Errorf("invalid input"))
			}
			d.initialPointer = v
			d.initialPoint = p
			return '.'
		}
		return v
	})
	if d.initialPointer == 0 {
		return fmt.Errorf("invalid input")
	}
	visited, _ := d.runBoard()
	fmt.Printf("Part 1: %d\n", len(visited))
	numLoops := 0
	for p := range visited {
		if d.board.Get(p) == '.' {
			d.board.Set(p, '#')
			_, loop := d.runBoard()
			if loop {
				numLoops++
			}
			d.board.Clear(p)
		}
	}
	fmt.Printf("Part 2: %d\n", numLoops)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
