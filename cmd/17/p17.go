package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ghjm/advent_utils"
)

type data struct {
	initState cpu
	program   []int
}

type cpu struct {
	regA int
	regB int
	regC int
	ip   int
}

func (d *data) newCPU() *cpu {
	c := &cpu{}
	c.regA = d.initState.regA
	c.regB = d.initState.regB
	c.regC = d.initState.regC
	return c
}

func (c *cpu) comboOperand(operand int) int {
	if operand >= 0 && operand <= 3 {
		return operand
	}
	if operand == 4 {
		return c.regA
	}
	if operand == 5 {
		return c.regB
	}
	if operand == 6 {
		return c.regC
	}
	return 0
}

func (c *cpu) runProgram(ctx context.Context, prog []int, out chan<- int) {
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if c.ip > len(prog)-2 {
				close(out)
				return
			}
			opcode := prog[c.ip]
			operand := prog[c.ip+1]
			c.ip += 2
			switch opcode {
			case 0: //adv
				operand = c.comboOperand(operand)
				divisor := 1
				for i := 0; i < operand; i++ {
					divisor *= 2
				}
				c.regA /= divisor
			case 1: //bxl
				c.regB ^= operand
			case 2: //bst
				operand = c.comboOperand(operand)
				c.regB = utils.Mod(operand, 8)
			case 3: //jnz
				if c.regA != 0 {
					c.ip = operand
				}
			case 4: //bxc
				c.regB = c.regB ^ c.regC
			case 5: //out
				operand = utils.Mod(c.comboOperand(operand), 8)
				out <- operand
			case 6: //bdv
				operand = c.comboOperand(operand)
				divisor := 1
				for i := 0; i < operand; i++ {
					divisor *= 2
				}
				c.regB = c.regA / divisor
			case 7: //cdv
				operand = c.comboOperand(operand)
				divisor := 1
				for i := 0; i < operand; i++ {
					divisor *= 2
				}
				c.regC = c.regA / divisor
			default:
				close(out)
				return
			}
		}
	}()
}

func (c *cpu) runProgramSync(prog []int) []int {
	ch := make(chan int)
	c.runProgram(context.Background(), prog, ch)
	var results []int
	for out := range ch {
		results = append(results, out)
	}
	return results
}

func (c *cpu) runOneStep(prog []int) int {
	ch := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	c.runProgram(ctx, prog, ch)
	result := <-ch
	cancel()
	go func() {
		for range ch {
		}
	}()
	return result
}

func outToStr(out []int) string {
	var tmp []string
	for _, v := range out {
		tmp = append(tmp, strconv.Itoa(v))
	}
	return strings.Join(tmp, ",")
}

// I have no idea how this works and just copied https://www.reddit.com/r/adventofcode/comments/1hg38ah/2024_day_17_solutions/m36hefc/
func (d *data) analyze(out []int, prevA int) int {
	if len(out) == 0 {
		return prevA
	}
	for a := range 1 << 10 {
		if a>>3 == prevA&127 {
			c := d.newCPU()
			c.regA = a
			r := c.runOneStep(d.program)
			if r == out[len(out)-1] {
				ret := d.analyze(out[:len(out)-1], (prevA<<3)|utils.Mod(a, 8))
				if ret >= 0 {
					return ret
				}
			}
		}
	}
	return -1
}

func run() error {
	d := data{
		initState: cpu{},
	}
	err := utils.OpenAndReadMultipleRegex("input17.txt", []utils.MultiRegex{
		{
			Regex: `^Register (\w): (\d+)$`,
			MatchFunc: func(s []string) error {
				switch s[1] {
				case "A":
					d.initState.regA = utils.MustAtoi(s[2])
				case "B":
					d.initState.regB = utils.MustAtoi(s[2])
				case "C":
					d.initState.regC = utils.MustAtoi(s[2])
				default:
					return fmt.Errorf("invalid register")
				}
				return nil
			},
		},
		{
			Regex: `^Program: (.+)$`,
			MatchFunc: func(s []string) error {
				p := strings.Split(s[1], ",")
				d.program = nil
				for _, pc := range p {
					d.program = append(d.program, utils.MustAtoi(pc))
				}
				return nil
			},
		},
	}, false)
	if err != nil {
		return err
	}
	{
		c := d.newCPU()
		out := c.runProgramSync(d.program)
		fmt.Printf("Part 1: %s\n", outToStr(out))
	}
	fmt.Printf("Part 2: %d\n", d.analyze(d.program, 0))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
