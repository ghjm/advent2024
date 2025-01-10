package main

import (
	"fmt"
	utils "github.com/ghjm/advent_utils"
	"github.com/ghjm/advent_utils/board"
	"github.com/ghjm/advent_utils/graph"
	"os"
)

type data struct {
	board    *board.StdBoard
	startPos utils.StdPoint
	endPos   utils.StdPoint
}

type state struct {
	pos       utils.StdPoint
	cheatUsed bool
}

func (d *data) getCheatCount(cheatSize int, threshold uint64) map[utils.StdRectangle]struct{} {
	startState := state{d.startPos, false}
	g := graph.Graph[state]{}
	g.CreateStateGraph(startState, func(s state) []graph.Edge[state] {
		var e []graph.Edge[state]
		for _, c := range d.board.Cardinals(s.pos, false) {
			if d.board.Get(c) == '.' {
				e = append(e, graph.Edge[state]{
					Dest: state{c, false},
					Cost: 1,
				})
			}
		}
		return e
	})
	gTmp := g.Copy()
	for n, edges := range gTmp.Nodes {
		n1 := state{n.pos, true}
		for _, e := range edges {
			n2 := state{e.Dest.pos, true}
			g.AddEdge(n1, n2, e.Cost)
		}
	}
	costs, _ := g.Dijkstra(startState)
	cheats := make(map[utils.StdRectangle]struct{})
	d.board.IterateBounds(func(u utils.Point[int]) bool {
		if d.board.Get(u) != '.' {
			return true
		}
		for dy := -cheatSize; dy <= cheatSize; dy++ {
			for dx := -cheatSize; dx <= cheatSize; dx++ {
				if dy == 0 && dx == 0 {
					continue
				}
				np := utils.StdPoint{X: u.X + dx, Y: u.Y + dy}
				if !d.board.Contains(np) {
					continue
				}
				if d.board.Get(np) != '.' {
					continue
				}
				md := uint64(u.ManhattanDistance(np))
				if md <= 1 || md >= uint64(cheatSize+1) {
					continue
				}
				cu := costs[state{u, false}]
				cnp := costs[state{np, false}]
				if cnp > cu+md+threshold-1 {
					cheats[utils.StdRectangle{
						P1: u,
						P2: np,
					}] = struct{}{}
				}
			}
		}
		return true
	})
	return cheats
}

func run() error {
	d := data{
		board: board.NewStdBoard(),
	}
	err := d.board.FromFile("input20.txt")
	const threshold = 100

	if err != nil {
		return err
	}
	d.board.Iterate(func(p utils.Point[int], v rune) bool {
		switch v {
		case 'S':
			d.startPos = p
			d.board.Set(p, '.')
		case 'E':
			d.endPos = p
			d.board.Set(p, '.')
		}
		return true
	})
	cheats := d.getCheatCount(2, 100)
	fmt.Printf("Part 1: %d\n", len(cheats))
	cheats = d.getCheatCount(20, 100)
	fmt.Printf("Part 2: %d\n", len(cheats))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
