package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	rules   map[int][]int
	updates [][]int
}

func (d *data) getRelevantRules(update []int) map[int][]int {
	updateSet := make(map[int]struct{})
	for _, v := range update {
		updateSet[v] = struct{}{}
	}
	relevantRules := make(map[int][]int)
	for k, v := range d.rules {
		_, ok := updateSet[k]
		if !ok {
			continue
		}
		var rr []int
		for _, r := range v {
			_, ok = updateSet[r]
			if ok {
				rr = append(rr, r)
			}
		}
		if len(rr) > 0 {
			relevantRules[k] = rr
		}
	}
	return relevantRules
}

func (d *data) checkValid(update []int) bool {
	relevantRules := d.getRelevantRules(update)
	seenSet := make(map[int]struct{})
	for _, v := range update {
		rr, ok := relevantRules[v]
		if ok {
			for _, r := range rr {
				_, ok = seenSet[r]
				if !ok {
					return false
				}
			}
		}
		seenSet[v] = struct{}{}
	}
	return true
}

func (d *data) reorder(update []int) []int {
	relevantRules := d.getRelevantRules(update)
	var result []int
	seenSet := make(map[int]struct{})
	for _, v := range update {
		_, ok := seenSet[v]
		if ok {
			continue
		}
		rr, ok := relevantRules[v]
		if ok {
			for _, r := range rr {
				_, ok = seenSet[r]
				if !ok {
					result = append(result, r)
					seenSet[r] = struct{}{}
				}
			}
		}
		result = append(result, v)
		seenSet[v] = struct{}{}
	}
	if reflect.DeepEqual(update, result) {
		return result
	} else {
		return d.reorder(result)
	}
}

func run() error {
	d := data{
		rules: make(map[int][]int),
	}
	readingUpdates := false
	err := utils.OpenAndReadLines("input5.txt", func(s string) error {
		if s == "" {
			readingUpdates = true
			return nil
		}
		if readingUpdates {
			var u []int
			for _, v := range strings.Split(s, ",") {
				u = append(u, utils.MustAtoi(v))
			}
			d.updates = append(d.updates, u)
		} else {
			r := strings.Split(s, "|")
			if len(r) != 2 {
				return fmt.Errorf("invalid input")
			}
			k, v := utils.MustAtoi(r[1]), utils.MustAtoi(r[0])
			_, ok := d.rules[k]
			if !ok {
				d.rules[k] = make([]int, 0)
			}
			d.rules[k] = append(d.rules[k], v)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sum := 0
	reorderSum := 0
	for _, u := range d.updates {
		if d.checkValid(u) {
			sum += u[len(u)/2]
		} else {
			r := d.reorder(u)
			reorderSum += r[len(r)/2]
		}
	}
	fmt.Printf("part 1: %d\n", sum)
	fmt.Printf("part 2: %d\n", reorderSum)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
