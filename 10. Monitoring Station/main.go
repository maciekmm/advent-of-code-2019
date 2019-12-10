package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

type boolGrid [][]bool

var grid boolGrid

type pair [2]int

func (bg boolGrid) At(co pair) bool {
	return bg[co[0]][co[1]]
}

func (bg pair) Add(vec pair) pair {
	return [2]int{bg[0] + vec[0], bg[1] + vec[1]}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func gridUnit(from pair, to pair) pair {
	vec := [2]int{to[0] - from[0], to[1] - from[1]}
	den := gcd(vec[0], vec[1])
	return [2]int{vec[0] / den, vec[1] / den}
}

func (bg boolGrid) Visible(station pair, planet pair) (bool, int) {
	step := gridUnit(station, planet)
	current := station.Add(step)
	covers := 0
	for current != planet {
		if bg.At(current) {
			covers++
		}
		current = current.Add(step)
	}
	return covers == 0, covers
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type relativePos struct {
	coord       pair
	cannonAngle float64
}

func main() {
	strgr, err := ioutil.ReadFile("input")
	checkErr(err)
	rows := strings.Split(string(strgr), "\n")
	cols := len(rows[0])
	grid = make([][]bool, len(rows))
	for i, row := range rows {
		grid[i] = make([]bool, cols)
		for j, cell := range row {
			if cell == '#' {
				grid[i][j] = true
			}
		}
	}
	max := 0
	var asteroids []relativePos
	for i, stationRow := range grid {
		for j, stationCell := range stationRow {
			if !stationCell {
				continue
			}
			station := pair{i, j}
			var asteroidsRelativePos = []relativePos{}
			reachable := 0
			for k, asteroidRow := range grid {
				for l, asteroidCell := range asteroidRow {
					if (k == i && j == l) || !asteroidCell {
						continue
					}
					visible, asteroidsBetween := grid.Visible(station, pair{k, l})
					if visible {
						reachable++
					}
					asteroidsRelativePos = append(asteroidsRelativePos, relativePos{
						cannonAngle: -(math.Atan2(float64(l-j), float64(k-i)) - math.Pi/2) + math.Pi*2*float64(asteroidsBetween),
						coord:       pair{k, l},
					})
				}
			}
			if reachable > max {
				max = reachable
				asteroids = asteroidsRelativePos
			}
		}
	}
	sort.Slice(asteroids, func(i, j int) bool {
		return asteroids[i].cannonAngle < asteroids[j].cannonAngle
	})
	fmt.Println("Part 1:", max)
	twohundreth := asteroids[199]
	fmt.Println("Part 2:", twohundreth.coord[1]*100+twohundreth.coord[0])
}
