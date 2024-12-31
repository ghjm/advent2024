package main

import (
	"fmt"
	utils "github.com/ghjm/advent_utils"
	"github.com/ghjm/advent_utils/board"
	"github.com/ghjm/advent_utils/graph"
	"math"
	"os"
)

var directions []utils.StdPoint = []utils.StdPoint{
	{X: 1, Y: 0},
	{X: 0, Y: 1},
	{X: -1, Y: 0},
	{X: 0, Y: -1},
}

type state struct {
	pos    utils.StdPoint
	facing int
}

type data struct {
	board        *board.StdBoard
	graph        *graph.Graph[state]
	initialState state
	terminalPos  utils.StdPoint
}

func run() error {
	d := data{
		board: board.NewStdBoard(),
		graph: &graph.Graph[state]{},
	}
	err := d.board.FromFile("input16.txt")
	if err != nil {
		return err
	}
	d.board.Iterate(func(p utils.Point[int], v rune) bool {
		switch v {
		case 'S':
			d.initialState = state{
				pos:    p,
				facing: 0,
			}
		case 'E':
			d.terminalPos = p
		}
		return true
	})
	d.graph.BuildStateGraph(d.initialState, func(s state) []graph.Edge[state] {
		var edges []graph.Edge[state]
		edges = append(edges, graph.Edge[state]{
			Dest: state{s.pos, utils.Mod(s.facing+1, 4)},
			Cost: 1000,
		})
		edges = append(edges, graph.Edge[state]{
			Dest: state{s.pos, utils.Mod(s.facing-1, 4)},
			Cost: 1000,
		})
		np := s.pos.Add(directions[s.facing])
		if d.board.Contains(np) {
			npv := d.board.Get(np)
			if npv == '.' || npv == 'S' || npv == 'E' {
				edges = append(edges, graph.Edge[state]{Dest: state{np, s.facing}, Cost: 1})
			}
		}
		return edges
	})
	cost, prev := d.graph.Dijkstra(d.initialState)
	var bestCost uint64 = math.MaxUint64
	var bestEnds []state
	for k, v := range cost {
		if k.pos == d.terminalPos {
			if v < bestCost {
				bestCost = v
				bestEnds = []state{k}
			} else if v == bestCost {
				bestEnds = append(bestEnds, k)
			}
		}
	}
	fmt.Printf("Part 1: %d\n", bestCost)

	visited := make(map[utils.StdPoint]struct{})
	open := bestEnds[:]
	for len(open) > 0 {
		s := open[0]
		open = open[1:]
		visited[s.pos] = struct{}{}
		pr, ok := prev[s]
		if ok {
			for _, p := range pr {
				open = append(open, p)
			}
		}
	}
	fmt.Printf("Part 2: %d\n", len(visited))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
