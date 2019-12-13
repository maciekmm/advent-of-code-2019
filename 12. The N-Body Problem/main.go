package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type vector3d struct {
	x int
	y int
	z int
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (v vector3d) absSum() int {
	return abs(v.x) + abs(v.y) + abs(v.z)
}

func (v vector3d) add(ot vector3d) vector3d {
	return vector3d{
		v.x + ot.x,
		v.y + ot.y,
		v.z + ot.z,
	}
}

func parseVector(vec string) (vector3d, error) {
	res := vector3d{}
	parts := strings.Split(strings.Trim(vec, "<>"), " ")
	for i, val := range []*int{&res.x, &res.y, &res.z} {
		num, err := strconv.Atoi(strings.Split(strings.Trim(parts[i], ","), "=")[1])
		if err != nil {
			return res, err
		}
		*val = num
	}
	return res, nil
}

func totalEnergy(position []vector3d, velocities []vector3d) int {
	energy := 0
	for i, pos := range position {
		energy += pos.absSum() * velocities[i].absSum()
	}
	return energy
}

func stateKey(dim int, positions []vector3d, velocities []vector3d) string {
	out := ""
	for i, pl := range positions {
		for j, val := range []*int{&pl.x, &pl.y, &pl.z} {
			if j == dim {
				out += strconv.Itoa(*val) + ";"
			}
		}
		for j, val := range []*int{&velocities[i].x, &velocities[i].y, &velocities[i].z} {
			if j == dim {
				out += strconv.Itoa(*val) + ";"
			}
		}
	}
	return out
}

func period(positions []vector3d) vector3d {
	visitedStates := [3]map[string]interface{}{}
	for i := 0; i < 3; i++ {
		visitedStates[i] = make(map[string]interface{})
	}
	velocities := make([]vector3d, len(positions))
	period := vector3d{}
	its := 0
	for period.x == 0 || period.y == 0 || period.z == 0 {
		its++
		for i := 0; i < 3; i++ {
			visitedStates[i][stateKey(i, positions, velocities)] = struct{}{}
		}
		step(positions, velocities)
		for i, per := range []*int{&period.x, &period.y, &period.z} {
			if _, ok := visitedStates[i][stateKey(i, positions, velocities)]; ok {
				if *per == 0 {
					*per = its
				}
			}
			visitedStates[i][stateKey(i, positions, velocities)] = struct{}{}
		}
	}
	return period
}

func step(positions []vector3d, velocities []vector3d) {
	for j, pos := range positions {
		for k, spos := range positions {
			if k <= j {
				continue
			}
			if pos.x < spos.x {
				velocities[j].x++
				velocities[k].x--
			} else if pos.x > spos.x {
				velocities[j].x--
				velocities[k].x++
			}
			if pos.y < spos.y {
				velocities[j].y++
				velocities[k].y--
			} else if pos.y > spos.y {
				velocities[j].y--
				velocities[k].y++
			}
			if pos.z < spos.z {
				velocities[j].z++
				velocities[k].z--
			} else if pos.z > spos.z {
				velocities[j].z--
				velocities[k].z++
			}
		}
		positions[j] = positions[j].add(velocities[j])
	}
}

func simulate(positions []vector3d, iterations int) ([]vector3d, []vector3d) {
	velocities := make([]vector3d, len(positions))
	for i := 0; i < iterations; i++ {
		step(positions, velocities)
	}
	return positions, velocities
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

func lcm(a, b int) int {
	den := gcd(a, b)
	return a / den * b
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	b, err := ioutil.ReadFile("input")
	checkErr(err)
	moonStrs := strings.Split(string(b), "\n")
	var positions []vector3d
	for _, moonStr := range moonStrs {
		pos, err := parseVector(moonStr)
		checkErr(err)
		positions = append(positions, pos)
	}
	part1 := make([]vector3d, len(positions))
	copy(part1, positions)
	fmt.Println("Part 1: ", totalEnergy(simulate(part1, 1000)))
	totalPeriod := period(positions)
	fmt.Println("Part 2: ", lcm(totalPeriod.x, lcm(totalPeriod.y, totalPeriod.z)))
}
