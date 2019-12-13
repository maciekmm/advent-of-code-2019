package intcode

import (
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
	OpBadj = 9
	OpHalt = 99
)

const (
	ModeAbsolute  = 0
	ModeImmediate = 1
	ModeRelative  = 2
)

type program map[int]int

var paramCount map[int]int = map[int]int{OpAdd: 3, OpMul: 3, OpIn: 1, OpOut: 1, OpJit: 2, OpJif: 2, OpLt: 3, OpEq: 3, OpHalt: 0, OpBadj: 1}

func (p program) Run(inp chan int, out chan int) {
	base := 0
	for i := 0; ; i++ {
		opCode := p[i] % 100
		mode := [3]int{
			(p[i] / 100) % 10,
			(p[i] / 100 / 10) % 10,
			(p[i] / 100 / 10 / 10) % 10,
		}
		var paramPointers [3]int
		var params [3]int
		for j := 0; j < paramCount[opCode]; j++ {
			mode := mode[j]
			if mode == ModeImmediate {
				paramPointers[j] = i + j + 1
			} else if mode == ModeAbsolute {
				paramPointers[j] = p[i+j+1]
			} else {
				paramPointers[j] = p[i+j+1] + base
			}
			params[j] = p[paramPointers[j]]
		}
		switch opCode {
		case OpAdd:
			p[paramPointers[2]] = params[0] + params[1]
		case OpMul:
			p[paramPointers[2]] = params[0] * params[1]
		case OpIn:
			val := <-inp
			p[paramPointers[0]] = val
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
				p[paramPointers[2]] = 1
			} else {
				p[paramPointers[2]] = 0
			}
		case OpEq:
			if params[0] == params[1] {
				p[paramPointers[2]] = 1
			} else {
				p[paramPointers[2]] = 0
			}
		case OpBadj:
			base += params[0]
		case OpHalt:
			close(out)
			return
		default:
		}
		i += paramCount[opCode]
	}
}

func (p program) Fork() program {
	prog := make(map[int]int)
	for k, v := range p {
		prog[k] = v
	}
	return prog
}

type Computer struct {
	Memory program
	Input  chan int
	Output chan int
}

func (c *Computer) Run() {
	c.Memory.Run(c.Input, c.Output)
}

func (c *Computer) Fork() *Computer {
	return &Computer{
		Memory: c.Memory.Fork(),
		Input:  make(chan int, 1),
		Output: make(chan int, 1),
	}
}

func LoadFile(in string) (*Computer, error) {
	f, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	cont, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	program := make(map[int]int)
	intCodes := strings.Split(string(cont), ",")
	for pos, intCodeStr := range intCodes {
		intCode, err := strconv.Atoi(strings.Trim(intCodeStr, "\n"))
		if err != nil {
			return nil, err
		}
		program[pos] = intCode
	}
	return &Computer{
		Memory: program,
		Input:  make(chan int, 1),
		Output: make(chan int, 1),
	}, nil
}
