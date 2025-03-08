package main

import (
	"fmt"
	"os"

	"github.com/ghjm/advent_utils"
)

type monkeySecret = uint64

type monkeyPrice struct {
	price  uint
	change int
}

type quad = [4]int

type data struct {
	monkeySecrets []monkeySecret
	monkeyPrices  [][]monkeyPrice
	quadsMap      map[quad]map[int]int
}

func mixPrune(v1, v2 monkeySecret) monkeySecret {
	r := v1 ^ v2
	return r % 16777216
}

func evolve(v monkeySecret) monkeySecret {
	r := mixPrune(v, v*64)
	r = mixPrune(r, r/32)
	r = mixPrune(r, r*2048)
	return r
}

func repeatEvolve(v monkeySecret, repeats uint) monkeySecret {
	r := v
	for i := uint(0); i < repeats; i++ {
		r = evolve(r)
	}
	return r
}

func (d *data) calcMonkeyPrices() {
	d.monkeyPrices = make([][]monkeyPrice, 0)
	for _, secret := range d.monkeySecrets {
		mp := make([]monkeyPrice, 0)
		prev := secret
		for i := 0; i < 2000; i++ {
			cur := evolve(prev)
			mp = append(mp, monkeyPrice{
				price:  uint(cur % 10),
				change: int(cur%10) - int(prev%10),
			})
			prev = cur
		}
		d.monkeyPrices = append(d.monkeyPrices, mp)
	}
}

func (d *data) calcQuadsMap() {
	d.quadsMap = make(map[quad]map[int]int)
	for i := range d.monkeyPrices {
		var rolling []int
		for j, p := range d.monkeyPrices[i] {
			rolling = append(rolling, p.change)
			for len(rolling) > 4 {
				rolling = rolling[1:]
			}
			if len(rolling) == 4 {
				q := quad{rolling[0], rolling[1], rolling[2], rolling[3]}
				_, ok := d.quadsMap[q]
				if !ok {
					d.quadsMap[q] = make(map[int]int)
				}
				_, ok = d.quadsMap[q][i]
				if !ok {
					d.quadsMap[q][i] = j
				}
			}
		}
	}
}

func (d *data) getTotalBananas(q quad) uint64 {
	var r uint64
	for i := range d.monkeyPrices {
		j, ok := d.quadsMap[q][i]
		if ok {
			r += uint64(d.monkeyPrices[i][j].price)
		}
	}
	return r
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input22.txt", func(s string) error {
		d.monkeySecrets = append(d.monkeySecrets, utils.MustAtoiU64(s))
		return nil
	})
	if err != nil {
		return err
	}
	sum := monkeySecret(0)
	for _, m := range d.monkeySecrets {
		sum += repeatEvolve(m, 2000)
	}
	fmt.Printf("Part 1: %d\n", sum)
	d.calcMonkeyPrices()
	d.calcQuadsMap()
	var maxBananas uint64
	for q := range d.quadsMap {
		b := d.getTotalBananas(q)
		if b > maxBananas {
			maxBananas = b
		}
	}
	fmt.Printf("Part 2: %d\n", maxBananas)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
