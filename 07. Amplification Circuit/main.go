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

//https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

type amplifier struct {
	Program Program
	in      chan int
	out     chan int
}

func runAmplifierCircuit(phases []int, program Program) (power int) {
	amplifiers := []amplifier{}
	prevOutput := make(chan int, 2)
	for i, phase := range phases {
		amp := amplifier{
			Program: program.Clone(),
			in:      prevOutput,
			out:     make(chan int, 2),
		}
		amplifiers = append(amplifiers, amp)
		amp.in <- phase
		if i == 0 {
			amp.in <- 0
		}
		prevOutput = amp.out
	}
	for _, amp := range amplifiers {
		go amp.Program.Run(amp.in, amp.out)
	}
	// feedback loop and power check
	for power = range prevOutput {
		amplifiers[0].in <- power
	}
	return
}

func maxThrusterPower(phases []int, programTemplate Program) int {
	perms := permutations(phases)
	max := 0
	for _, perm := range perms {
		power := runAmplifierCircuit(perm, programTemplate.Clone())
		if power > max {
			max = power
		}
	}
	return max
}

func main() {
	data, err := loadData()
	checkErr(err)
	phases := []int{0, 1, 2, 3, 4}
	fmt.Println("Part 1:", maxThrusterPower(phases, data))
	phases = []int{5, 6, 7, 8, 9}
	fmt.Println("Part 2:", maxThrusterPower(phases, data))
}
