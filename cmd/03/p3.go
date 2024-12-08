package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type mulinstr struct {
	a int
	b int
}

type data struct {
	muls        []mulinstr
	enabledMuls []mulinstr
}

func run() error {
	d := data{}
	re1 := regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	enabled := true
	err := utils.OpenAndReadLines("input3.txt", func(s string) error {
		m := re1.FindAllStringSubmatch(s, -1)
		for _, sm := range m {
			if strings.HasPrefix(sm[0], "mul") {
				mi := mulinstr{
					a: utils.MustAtoi(sm[1]),
					b: utils.MustAtoi(sm[2]),
				}
				d.muls = append(d.muls, mi)
				if enabled {
					d.enabledMuls = append(d.enabledMuls, mi)
				}
			} else if strings.HasPrefix(sm[0], "don't") {
				enabled = false
			} else if strings.HasPrefix(sm[0], "do") {
				enabled = true
			} else {
				return fmt.Errorf("impossible match")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	sum := 0
	for _, m := range d.muls {
		sum += m.a * m.b
	}
	fmt.Printf("Part 1: %d\n", sum)
	sum = 0
	for _, m := range d.enabledMuls {
		sum += m.a * m.b
	}
	fmt.Printf("Part 2: %d\n", sum)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
