package main

import (
	"fmt"
	"github.com/ghjm/advent2024/pkg/utils"
	"os"
	"sort"
	"strings"
)

type data struct {
	list1 []int
	list2 []int
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input1.txt", func(s string) error {
		if s == "" {
			return nil
		}
		var s2 string
		for s != s2 {
			s2 = s
			s = strings.Replace(s2, "  ", " ", -1)
		}
		l := strings.Split(s, " ")
		if len(l) != 2 {
			return fmt.Errorf("invalid input")
		}
		d.list1 = append(d.list1, utils.MustAtoi(strings.TrimSpace(l[0])))
		d.list2 = append(d.list2, utils.MustAtoi(strings.TrimSpace(l[1])))
		return nil
	})
	if err != nil {
		return err
	}
	if len(d.list1) != len(d.list2) {
		return fmt.Errorf("invalid input")
	}
	sort.Ints(d.list1)
	sort.Ints(d.list2)
	sum := 0
	for i := range d.list1 {
		sum += utils.Abs(d.list1[i] - d.list2[i])
	}
	fmt.Printf("Part 1: %d\n", sum)
	list2vals := make(map[int]int)
	for _, v := range d.list2 {
		prev, ok := list2vals[v]
		if ok {
			list2vals[v] = prev + 1
		} else {
			list2vals[v] = 1
		}
	}
	similarity := 0
	for _, v := range d.list1 {
		app, ok := list2vals[v]
		if !ok {
			app = 0
		}
		similarity += app * v
	}
	fmt.Printf("Part 2: %d\n", similarity)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
