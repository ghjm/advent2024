package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ghjm/advent_utils"
)

type stoneMap map[uint64]uint64

type data struct {
	stones stoneMap
}

func (sm *stoneMap) addStones(id, count uint64) {
	_, ok := (*sm)[id]
	if !ok {
		(*sm)[id] = count
	} else {
		(*sm)[id] += count
	}
}

func (sm *stoneMap) copy() stoneMap {
	result := make(stoneMap)
	for k, v := range *sm {
		result[k] = v
	}
	return result
}

func (sm *stoneMap) count() uint64 {
	var result uint64
	for _, v := range *sm {
		result += v
	}
	return result
}

func transform(stones stoneMap) stoneMap {
	results := make(stoneMap)
	for st, count := range stones {
		stStr := strconv.FormatUint(st, 10)
		if st == 0 {
			results.addStones(1, count)
		} else if len(stStr)%2 == 0 {
			results.addStones(utils.MustAtoiU64(stStr[:len(stStr)/2]), count)
			results.addStones(utils.MustAtoiU64(stStr[len(stStr)/2:]), count)
		} else {
			results.addStones(st*2024, count)
		}
	}
	return results
}

func run() error {
	d := data{
		stones: make(stoneMap),
	}
	s, err := utils.OpenAndReadAll("input11.txt")
	if err != nil {
		return err
	}
	for _, st := range strings.Split(string(s), " ") {
		d.stones.addStones(utils.MustAtoiU64(st), 1)
	}
	st := d.stones.copy()
	for i := 0; i < 25; i++ {
		st = transform(st)
	}
	fmt.Printf("Part 1: %d\n", st.count())
	for i := 0; i < 50; i++ {
		st = transform(st)
	}
	fmt.Printf("Part 2: %d\n", st.count())
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
