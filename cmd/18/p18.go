package main

import (
	"fmt"
	"github.com/ghjm/advent_utils/board"
	"github.com/ghjm/advent_utils/graph"
	"os"
	"strings"

	"github.com/ghjm/advent_utils"
)

type data struct {
	board        *board.StdBoard
	fallingBytes []utils.StdPoint
}

func boardToGraph(b *board.StdBoard) graph.Graph[utils.StdPoint] {
	g := graph.Graph[utils.StdPoint]{}
	g.CreateStateGraph(utils.StdPoint{}, func(point utils.StdPoint) []graph.Edge[utils.StdPoint] {
		var e []graph.Edge[utils.StdPoint]
		for _, np := range b.Cardinals(point, false) {
			if b.Get(np) != '#' {
				e = append(e, graph.Edge[utils.StdPoint]{
					Dest: np,
					Cost: 1,
				})
			}
		}
		return e
	})
	return g
}

func run() error {
	d := data{
		board: board.NewStdBoard(board.WithBounds[int, rune](utils.Rectangle[int]{
			P2: utils.Point[int]{
				X: 70,
				Y: 70,
			},
		})),
	}
	err := utils.OpenAndReadLines("input18.txt", func(s string) error {
		fb := strings.Split(s, ",")
		if len(fb) != 2 {
			return fmt.Errorf("invalid input")
		}
		d.fallingBytes = append(d.fallingBytes, utils.StdPoint{
			X: utils.MustAtoi(fb[0]),
			Y: utils.MustAtoi(fb[1]),
		})
		return nil
	})
	if err != nil {
		return err
	}
	bp1 := d.board.Copy()
	for i := 0; i < 1024; i++ {
		bp1.Set(d.fallingBytes[i], '#')
	}
	start := utils.StdPoint{}
	finish := utils.StdPoint{X: 70, Y: 70}
	{
		g := boardToGraph(bp1)
		cost, _ := g.Dijkstra(start)
		c, ok := cost[finish]
		if !ok {
			return fmt.Errorf("no solution")
		}
		fmt.Printf("Part 1: %d\n", c)
	}
	for i := 1024; true; i++ {
		dp := d.fallingBytes[i]
		bp1.Set(dp, '#')
		g := boardToGraph(bp1)
		cost, _ := g.Dijkstra(start)
		_, ok := cost[finish]
		if !ok {
			fmt.Printf("Part 2: %d,%d", dp.X, dp.Y)
			break
		}
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
