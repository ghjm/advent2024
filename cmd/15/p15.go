package main

import (
	"errors"
	"fmt"
	"github.com/ghjm/advent_utils/board"
	"os"
	"strings"

	"github.com/ghjm/advent_utils"
)

type data struct {
	originalBoard *board.StdBoard
	board         *board.StdBoard
	moves         string
	robotPos      utils.StdPoint
}

func (d *data) tryMoveBox(p utils.StdPoint, dir utils.StdPoint) bool {
	dp := p.Add(dir)
	if !d.board.Contains(dp) {
		return false
	}
	dv := d.board.Get(dp)
	if dv == 'O' {
		if d.tryMoveBox(dp, dir) {
			dv = '.'
		} else {
			return false
		}
	}
	if dv == '.' {
		d.board.Set(dp, 'O')
		d.board.Set(p, '.')
		return true
	}
	return false
}

var directions = map[rune]utils.StdPoint{
	'<': utils.StdPoint{X: -1, Y: 0},
	'>': utils.StdPoint{X: 1, Y: 0},
	'^': utils.StdPoint{X: 0, Y: -1},
	'v': utils.StdPoint{X: 0, Y: 1},
}

func (d *data) executeMoves() {
	for _, m := range d.moves {
		dir, ok := directions[m]
		if !ok {
			panic("invalid move")
		}
		dp := d.robotPos.Add(dir)
		if !d.board.Contains(dp) {
			continue
		}
		dv := d.board.Get(dp)
		if dv == 'O' {
			if !d.tryMoveBox(dp, dir) {
				continue
			}
		} else if dv != '.' {
			continue
		}
		d.robotPos = dp
	}
}

func (d *data) getSumGPS() int {
	sum := 0
	d.board.Iterate(func(p utils.Point[int], v rune) bool {
		if v == 'O' {
			sum += (p.Y * 100) + p.X
		}
		return true
	})
	return sum
}

type box struct {
	pos utils.StdPoint
}

type dataP2 struct {
	board    *board.StdBoard
	boxes    []box
	moves    string
	robotPos utils.StdPoint
}

func (d *dataP2) findBox(p utils.StdPoint) int {
	for i, b := range d.boxes {
		if ((p.X == b.pos.X) || (p.X == b.pos.X+1)) && b.pos.Y == p.Y {
			return i
		}
	}
	return -1
}

func (d *dataP2) checkMoveBox(bx int, dir utils.StdPoint) (bool, map[int]struct{}) {
	var checkPoints []utils.StdPoint
	bp := d.boxes[bx].pos
	switch {
	case dir.X == -1 && dir.Y == 0:
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X - 1, Y: bp.Y})
	case dir.X == 1 && dir.Y == 0:
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X + 2, Y: bp.Y})
	case dir.X == 0 && dir.Y == -1:
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X, Y: bp.Y - 1})
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X + 1, Y: bp.Y - 1})
	case dir.X == 0 && dir.Y == 1:
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X, Y: bp.Y + 1})
		checkPoints = append(checkPoints, utils.StdPoint{X: bp.X + 1, Y: bp.Y + 1})
	default:
		panic("unsupported direction")
	}
	usedBoxes := make(map[int]struct{})
	usedBoxes[bx] = struct{}{}
	for _, cp := range checkPoints {
		if !d.board.Contains(cp) {
			return false, nil
		}
		if d.board.Get(cp) != '.' {
			return false, nil
		}
		cbx := d.findBox(cp)
		if cbx != -1 {
			ok, boxes := d.checkMoveBox(cbx, dir)
			if !ok {
				return false, nil
			}
			for b := range boxes {
				usedBoxes[b] = struct{}{}
			}
		}
	}
	return true, usedBoxes
}

func (d *dataP2) tryMoveBox(bx int, dir utils.StdPoint) bool {
	ok, boxes := d.checkMoveBox(bx, dir)
	if !ok {
		return false
	}
	for b := range boxes {
		d.boxes[b].pos = d.boxes[b].pos.Add(dir)
	}
	return true
}

func (d *dataP2) executeMoves() {
	for _, m := range d.moves {
		dir, ok := directions[m]
		if !ok {
			panic("invalid move")
		}
		dp := d.robotPos.Add(dir)
		if !d.board.Contains(dp) || d.board.Get(dp) != '.' {
			continue
		}
		bx := d.findBox(dp)
		if bx != -1 {
			if !d.tryMoveBox(bx, dir) {
				continue
			}
		}
		d.robotPos = dp
	}
}

func (d *data) convertToP2() *dataP2 {
	dp2 := &dataP2{
		board: board.NewStdBoard(board.WithStorage(&board.FlatBoard{})),
		moves: d.moves,
	}
	var newBoard []string
	for y := range d.originalBoard.Bounds().Height() {
		row := ""
		for x := range d.originalBoard.Bounds().Width() {
			switch d.originalBoard.Get(utils.StdPoint{X: x, Y: y}) {
			case '.':
				row += ".."
			case '#':
				row += "##"
			case 'O':
				row += ".."
				dp2.boxes = append(dp2.boxes, box{
					pos: utils.StdPoint{
						X: x * 2,
						Y: y,
					},
				})
			case '@':
				row += ".."
				dp2.robotPos = utils.StdPoint{
					X: x * 2,
					Y: y,
				}
			}
		}
		newBoard = append(newBoard, row)
	}
	err := dp2.board.FromStrings(newBoard)
	if err != nil {
		panic(err)
	}
	return dp2
}

func (d *dataP2) getSumGPS() int {
	sum := 0
	for _, b := range d.boxes {
		sum += (b.pos.Y * 100) + b.pos.X
	}
	return sum
}

func (d *dataP2) print() {
	var bd [][]rune
	for _, str := range d.board.Format() {
		bd = append(bd, []rune(str))
	}
	bd[d.robotPos.Y][d.robotPos.X] = '@'
	for _, b := range d.boxes {
		bd[b.pos.Y][b.pos.X] = '['
		bd[b.pos.Y][b.pos.X+1] = ']'
	}
	for _, row := range bd {
		for _, ch := range row {
			fmt.Print(string(ch))
		}
		fmt.Println("")
	}
}

func run() error {
	d := data{
		board: board.NewStdBoard(board.WithStorage(&board.FlatBoard{})),
	}
	f, err := utils.OpenAndReadAll("input15.txt")
	if err != nil {
		return err
	}
	fs := strings.Split(string(f), "\n\n")
	if len(fs) != 2 {
		return errors.New("invalid input")
	}
	err = d.board.FromStrings(strings.Split(fs[0], "\n"))
	if err != nil {
		return err
	}
	d.originalBoard = d.board.Copy()
	d.moves = strings.Replace(fs[1], "\n", "", -1)
	d.board.Iterate(func(p utils.Point[int], v rune) bool {
		if v == '@' {
			d.board.Set(p, '.')
			d.robotPos = p
			return false
		}
		return true
	})
	d.executeMoves()
	fmt.Printf("Part 1: %d\n", d.getSumGPS())
	d2 := d.convertToP2()
	d2.executeMoves()
	fmt.Printf("Part 2: %d\n", d2.getSumGPS())
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
