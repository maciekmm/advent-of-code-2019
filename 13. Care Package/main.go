package main

import "github.com/maciekmm/aoc-2019/intcode"

import "fmt"

type tileType int

const (
	tileEmpty tileType = 0
	tileWall  tileType = 1
	tileBlock tileType = 2
	tilePadle tileType = 3
	tileBall  tileType = 4
)

type position struct {
	x, y int
}

type game struct {
	cmp    *intcode.Computer
	ball   position
	padle  position
	score  int
	blocks map[position]interface{}
}

func (game *game) Run(input func(g *game) int) {
	go game.cmp.Run()
	for {
		x, ok := <-game.cmp.Output
		if !ok {
			break
		}
		y, ok := <-game.cmp.Output
		if !ok {
			break
		}
		tile, ok := <-game.cmp.Output
		if !ok {
			break
		}
		pos := position{x, y}
		if x == -1 && y == 0 {
			game.score = tile
			continue
		}
		switch tileType(tile) {
		case tileEmpty:
			delete(game.blocks, pos)
		case tilePadle:
			game.padle = pos
		case tileBlock:
			game.blocks[pos] = struct{}{}
		case tileBall:
			game.ball = pos
			if input != nil {
				game.cmp.Input <- input(game)
			}
		}
	}
}

func main() {
	cmp, err := intcode.LoadFile("input")
	if err != nil {
		panic(err)
	}
	g := &game{
		cmp:    cmp.Fork(),
		blocks: make(map[position]interface{}),
	}
	g.Run(nil)
	fmt.Println("Part 1: ", len(g.blocks))

	cmp.Memory[0] = 2
	game2 := &game{
		cmp:    cmp.Fork(),
		blocks: make(map[position]interface{}),
	}
	game2.Run(func(g *game) int {
		if g.padle.x > g.ball.x {
			return -1
		} else if g.padle.x < g.ball.x {
			return 1
		}
		return 0
	})
	fmt.Println("Part 2: ", game2.score)
}
