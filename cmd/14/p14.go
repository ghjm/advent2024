package main

import (
	"fmt"
	"github.com/ghjm/advent_utils/board"
	"os"

	"github.com/ghjm/advent_utils"
)

type robot struct {
	pos utils.Point[int64]
	vel utils.Point[int64]
}

type data struct {
	robots      []robot
	boardWidth  int
	boardHeight int
}

func (d *data) step(t int64) *data {
	var result data
	result.boardWidth = d.boardWidth
	result.boardHeight = d.boardHeight
	for _, rob := range d.robots {
		result.robots = append(result.robots, rob.futurePos(t).modBoard(int64(d.boardWidth), int64(d.boardHeight)))
	}
	return &result
}

func (r robot) futurePos(t int64) robot {
	return robot{
		pos: utils.Point[int64]{
			X: r.pos.X + r.vel.X*t,
			Y: r.pos.Y + r.vel.Y*t,
		},
		vel: r.vel,
	}
}

func (r robot) modBoard(width, height int64) robot {
	rp := robot{
		pos: utils.Point[int64]{
			X: r.pos.X % width,
			Y: r.pos.Y % height,
		},
		vel: r.vel,
	}
	for rp.pos.X < 0 {
		rp.pos.X += width
	}
	for rp.pos.Y < 0 {
		rp.pos.Y += height
	}
	return rp
}

func run() error {
	d := data{}
	//d.boardWidth, d.boardHeight = 11, 7
	d.boardWidth, d.boardHeight = 101, 103
	m, err := utils.OpenAndReadRegex("input14.txt", `p=(\d+),(\d+) v=(-?\d+),(-?\d+)`, true)
	if err != nil {
		return err
	}
	for _, v := range m {
		r := robot{
			pos: utils.Point[int64]{
				X: utils.MustAtoi64(v[1]),
				Y: utils.MustAtoi64(v[2]),
			},
			vel: utils.Point[int64]{
				X: utils.MustAtoi64(v[3]),
				Y: utils.MustAtoi64(v[4]),
			},
		}
		d.robots = append(d.robots, r)
	}
	var fR []robot
	for _, r := range d.robots {
		fR = append(fR, r.futurePos(100).modBoard(int64(d.boardWidth), int64(d.boardHeight)))
	}
	fb := board.FlatBoard{}
	fb.Allocate(d.boardWidth, d.boardHeight, '0')
	b := board.NewStdBoard(board.WithStorage(&fb), board.WithBounds[int, rune](fb.GetBounds()), board.WithEmptyVal[int, rune]('0'))
	b.IterateBounds(func(u utils.Point[int]) bool {
		b.Set(u, '0')
		return true
	})
	for _, r := range fR {
		rp := utils.StdPoint{X: int(r.pos.X), Y: int(r.pos.Y)}
		b.Set(rp, b.Get(rp)+1)
	}
	var part1 int64 = 1
	for qY := range 2 {
		for qX := range 2 {
			var sum int64
			for y := range d.boardHeight / 2 {
				for x := range d.boardWidth / 2 {
					v := b.Get(utils.StdPoint{
						X: qX*(d.boardWidth/2+1) + x,
						Y: qY*(d.boardHeight/2+1) + y,
					})
					sum += utils.MustAtoi64(string(v))
				}
			}
			part1 *= sum
		}
	}
	fmt.Printf("Part 1: %d\n", part1)
	fD := d.step(1)
	steps := 1
	for {
		fb := board.FlatBoard{}
		fb.Allocate(d.boardWidth, d.boardHeight, '.')
		b := board.NewStdBoard(board.WithStorage(&fb), board.WithBounds[int, rune](fb.GetBounds()), board.WithEmptyVal[int, rune]('.'))
		for _, r := range fD.robots {
			rp := utils.StdPoint{X: int(r.pos.X), Y: int(r.pos.Y)}
			b.Set(rp, '#')
		}
		reg := b.FindRegions(false)
		maxReg := 0
		for _, r := range reg {
			if len(r) > maxReg {
				maxReg = len(r)
			}
		}
		if maxReg > 200 {
			//b.Print()
			fmt.Printf("Part 2: %d\n", steps)
			break
		}
		fD = fD.step(1)
		steps++
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
