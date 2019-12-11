package main

import (
	"fmt"

	"github.com/maciekmm/aoc-2019/intcode"
)

type pair [2]int

func (p pair) Add(o pair) pair {
	return pair{p[0] + o[0], p[1] + o[1]}
}

const (
	colorBlack = 0
	colorWhite = 1
)

const (
	rotLeft  = 0
	rotRight = 1
)

func run(color int) map[pair]int {
	visitedTiles := make(map[pair]int)
	facings := [][2]int{[2]int{0, 1}, [2]int{1, 0}, [2]int{0, -1}, [2]int{-1, 0}}
	cmp, err := intcode.LoadFile("input")
	if err != nil {
		panic(err)
	}
	go cmp.Run()
	cmp.Input <- color
	pos := pair{0, 0}
	facing := 0
	for {
		select {
		case v, ok := <-cmp.Output:
			if !ok {
				return visitedTiles
			}
			visitedTiles[pos] = v
			if rot, ok := <-cmp.Output; ok {
				if rot == rotLeft {
					facing = (facing - 1 + len(facings)) % len(facings)
				} else {
					facing = (facing + 1) % len(facings)
				}
				pos = pos.Add(facings[facing])
				if v, ok := visitedTiles[pos]; ok {
					cmp.Input <- v
				} else {
					cmp.Input <- colorBlack
				}
			} else {
				return visitedTiles
			}
		}
	}
}

func draw(pl map[pair]int) {
	var f, t pair
	p := false
	for k := range pl {
		if !p || k[0] < f[0] {
			f[0] = k[0]
		}
		if !p || k[1] < f[1] {
			f[1] = k[1]
		}
		if !p || k[0] > t[0] {
			t[0] = k[0]
		}
		if !p || k[1] > t[1] {
			t[1] = k[1]
		}
		p = true
	}
	for y := t[1]; y >= f[1]; y-- {
		for x := f[0]; x < t[0]; x++ {
			color := colorBlack
			if k, ok := pl[pair{x, y}]; ok {
				color = k
			}
			if color == colorBlack {
				fmt.Printf("\u001b[30m█")
			} else {
				fmt.Printf("\u001b[37m█")
			}
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("Part 1: ", len(run(colorBlack)))
	draw(run(colorWhite))
}
