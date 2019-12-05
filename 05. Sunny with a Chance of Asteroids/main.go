package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	OpAdd  = 1
	OpMul  = 2
	OpIn   = 3
	OpOut  = 4
	OpJit  = 5
	OpJif  = 6
	OpLt   = 7
	OpEq   = 8
	OpHalt = 99
)

type Program []int

var paramCount map[int]int = map[int]int{OpAdd: 3, OpMul: 3, OpIn: 1, OpOut: 1, OpJit: 2, OpJif: 2, OpLt: 3, OpEq: 3, OpHalt: 0}

func (p Program) Run(inp chan int, out chan int) {
	for i := 0; i < len(p); i++ {
		opCode := p[i] % 100
		immediate := [3]bool{
			(p[i]/100)%10 == 1,
			(p[i]/100/10)%10 == 1,
			(p[i]/100/10/10)%10 == 1,
		}
		var params [3]int
		for j := 0; j < paramCount[opCode]; j++ {
			imm := immediate[j]
			if imm {
				params[j] = p[i+j+1]
			} else {
				params[j] = p[p[i+j+1]]
			}
		}
		switch opCode {
		case OpAdd:
			p[p[i+3]] = params[0] + params[1]
		case OpMul:
			p[p[i+3]] = params[0] * params[1]
		case OpIn:
			val := <-inp
			p[p[i+1]] = val
		case OpOut:
			out <- params[0]
		case OpJit:
			if params[0] != 0 {
				i = params[1] - 1
				continue
			}
		case OpJif:
			if params[0] == 0 {
				i = params[1] - 1
				continue
			}
		case OpLt:
			if params[0] < params[1] {
				p[p[i+3]] = 1
			} else {
				p[p[i+3]] = 0
			}
		case OpEq:
			if params[0] == params[1] {
				p[p[i+3]] = 1
			} else {
				p[p[i+3]] = 0
			}
		case OpHalt:
			close(out)
			return
		default:
		}
		i += paramCount[opCode]
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

func mustAtoi(str string) int {
	d, err := strconv.Atoi(str)
	checkErr(err)
	return d
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
		intCode, err := strconv.Atoi(strings.Trim(intCodeStr, "\n"))
		if err != nil {
			return nil, err
		}
		program[pos] = intCode
	}
	return Program(program), nil
}

func runWithInput(pro Program, input int) int {
	out := make(chan int)
	in := make(chan int, 1)
	in <- input
	go pro.Run(in, out)
	for val := range out {
		if val != 0 {
			return val
		}
	}
	return 0
}

func main() {
	data, err := loadData()
	checkErr(err)
	part1 := data.Clone()
	fmt.Printf("Part 1: %d\n", runWithInput(part1, 1))
	part2 := data.Clone()
	fmt.Printf("Part 2: %d\n", runWithInput(part2, 5))
}
