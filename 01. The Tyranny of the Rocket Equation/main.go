package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func calculateFuel(fuel int) int {
	if fuel <= 0 {
		return 0
	}
	fuel = fuel/3 - 2
	return fuel + calculateFuel(fuel)
}

func main() {
	part1 := 0
	part2 := 0
	f, err := os.Open("input")
	checkErr(err)
	cont, err := ioutil.ReadAll(f)
	checkErr(err)
	lines := strings.Split(string(cont), "\n")
	for _, spacecraft := range lines {
		if spacecraft == "" {
			break
		}
		mass, err := strconv.Atoi(spacecraft)
		checkErr(err)
		part1 += mass/3 - 2
		part2 += calculateFuel(mass)
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}
