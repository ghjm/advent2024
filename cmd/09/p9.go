package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ghjm/advent_utils"
)

func loadDisk(diskMap string) []int {
	var diskSectors []int
	for i, ch := range diskMap {
		size := utils.MustAtoi(string(ch))
		for range size {
			if i%2 == 0 {
				diskSectors = append(diskSectors, i/2)
			} else {
				diskSectors = append(diskSectors, -1)
			}

		}
	}
	return diskSectors
}

func defrag1(diskSectors []int) {
	freePtr := 0
	endPtr := len(diskSectors) - 1
	for endPtr > freePtr {
		for diskSectors[endPtr] == -1 {
			endPtr--
			if endPtr <= 0 {
				return
			}
		}
		for diskSectors[freePtr] != -1 {
			freePtr++
			if freePtr > endPtr {
				return
			}
		}
		diskSectors[freePtr] = diskSectors[endPtr]
		diskSectors[endPtr] = -1
	}
}

func findFile(diskSectors []int, id int) (start int, size int) {
	start = -1
	size = 0
	for i := range diskSectors {
		sector := diskSectors[i]
		if sector == id {
			if start == -1 {
				start = i
			}
			size++
		} else if start >= 0 && sector != id {
			return
		}
	}
	return
}

func findFree(diskSectors []int, reqSize int) int {
	start := -1
	size := 0
	for i, sector := range diskSectors {
		if sector == -1 {
			if start == -1 {
				start = i
				size = 0
			}
			size++
			if size >= reqSize {
				return start
			}
		} else {
			start = -1
			size = 0
		}
	}
	return -1
}

func defrag2(diskSectors []int) {
	topFile := slices.Max(diskSectors)
	for id := topFile; id >= 0; id-- {
		start, size := findFile(diskSectors, id)
		if start == -1 {
			panic(fmt.Sprintf("didn't find file %d", id))
		}
		free := findFree(diskSectors, size)
		if free == -1 || free > start {
			continue
		}
		for i := 0; i < size; i++ {
			diskSectors[free+i] = id
			diskSectors[start+i] = -1
		}
	}
}

func checksum(diskSectors []int) uint64 {
	var sum uint64
	for i, sector := range diskSectors {
		if sector >= 0 {
			sum += uint64(sector * i)
		}
	}
	return sum
}

func run() error {
	compact, err := utils.OpenAndReadAll("input9.txt")
	if err != nil {
		return err
	}

	ds := loadDisk(strings.TrimSpace(string(compact)))
	defrag1(ds)
	fmt.Printf("Part 1: %d\n", checksum(ds))

	ds = loadDisk(strings.TrimSpace(string(compact)))
	defrag2(ds)
	fmt.Printf("Part 2: %d\n", checksum(ds))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
