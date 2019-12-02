package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Program []int

func (p Program) Run() {
	for i := 0; i < len(p); i += 4 {
		switch p[i] {
		case 1:
			p[p[i+3]] = p[p[i+1]] + p[p[i+2]]
		case 2:
			p[p[i+3]] = p[p[i+1]] * p[p[i+2]]
		case 99:
			return
		default:
		}
	}
}

func (p Program) Clone() Program {
	prog := make([]int, len(p))
	copy(prog, p)
	return prog
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadData() (Program, error) {
	f, err := os.Open("input")
	if err != nil {
		return nil, err
	}
	cont, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	intCodes := strings.Split(string(cont), ",")
	program := make([]int, len(intCodes))
	for pos, intCodeStr := range intCodes {
		intCode, err := strconv.Atoi(intCodeStr)
		if err != nil {
			return nil, err
		}
		program[pos] = intCode
	}
	return Program(program), nil
}

func main() {
	data, err := loadData()
	part1 := data.Clone()
	checkErr(err)
	part1[1] = 12
	part1[2] = 2
	part1.Run()
	fmt.Printf("Part 1: %d\n", part1[0])
	//part 2
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			prog := data.Clone()
			prog[1] = noun
			prog[2] = verb
			prog.Run()
			if prog[0] == 19690720 {
				fmt.Printf("Part 2: %d\n", 100*noun+verb)
				return
			}
		}
	}
	fmt.Println("No value found")
}
