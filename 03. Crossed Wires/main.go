package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

var errNotIntersecting error = errors.New("no intersection")

type intersection struct {
	Coord coord
	Steps int
}

type coord struct {
	x int
	y int
}

func (c coord) Distance(o coord) int {
	return int(math.Abs(float64(c.x-o.x)) + math.Abs(float64(c.y-o.y)))
}

type cable []segment

func (c cable) Intersections(o cable) []intersection {
	var intersections []intersection
	firstPathSteps := 0
	firstPathPrevSeg := coord{0, 0}
	for _, cSeg := range c {
		firstPathSteps += cSeg.Start.Distance(firstPathPrevSeg)
		secondPathSteps := 0
		secondPathPrevSeg := coord{0, 0}
		for _, oSeg := range o {
			secondPathSteps += oSeg.Start.Distance(secondPathPrevSeg)
			in, err := cSeg.Intersection(oSeg)
			if err == nil {
				inter := intersection{*in, secondPathSteps + firstPathSteps + in.Distance(oSeg.Start) + in.Distance(cSeg.Start)}
				intersections = append(intersections, inter)
			}
			secondPathPrevSeg = oSeg.Start
		}
		firstPathPrevSeg = cSeg.Start
	}
	return intersections
}

type segment struct {
	Start coord
	End   coord
}

func (s segment) Vertical() bool {
	return s.Start.x == s.End.x
}

func (s segment) sort(byX bool) segment {
	if byX && s.Start.x > s.End.x {
		return segment{s.End, s.Start}
	} else if !byX && s.Start.y > s.End.y {
		return segment{s.End, s.Start}
	}
	return s
}

func (s segment) Intersection(o segment) (*coord, error) {
	// X axis
	oxed := o.sort(true)
	sxed := s.sort(true)
	if (o.Vertical() && (o.Start.x < sxed.Start.x || o.Start.x > sxed.End.x)) ||
		(s.Vertical() && (s.Start.x < oxed.Start.x || s.Start.x > oxed.End.x)) {
		return nil, errNotIntersecting
	}

	// Y Axis
	oxed = o.sort(false)
	sxed = s.sort(false)
	if (!o.Vertical() && (o.Start.y < sxed.Start.y || o.Start.y > sxed.End.y)) ||
		(!s.Vertical() && (s.Start.y < oxed.Start.y || s.Start.y > oxed.End.y)) {
		return nil, errNotIntersecting
	}

	if o.Vertical() {
		return &coord{o.Start.x, s.Start.y}, nil
	}
	return &coord{s.Start.x, o.Start.y}, nil
}

func assembleCircuit(repr string) (cable, error) {
	var cable []segment
	var seg segment
	prev := coord{0, 0}
	instrs := strings.Split(repr, ",")
	for _, instr := range instrs {
		dist, err := strconv.Atoi(instr[1:])
		if err != nil {
			return nil, err
		}
		switch instr[0] {
		case 'U':
			seg = segment{Start: prev, End: coord{prev.x, prev.y + dist}}
		case 'D':
			seg = segment{Start: prev, End: coord{prev.x, prev.y - dist}}
		case 'L':
			seg = segment{Start: prev, End: coord{prev.x - dist, prev.y}}
		case 'R':
			seg = segment{Start: prev, End: coord{prev.x + dist, prev.y}}
		}

		prev = seg.End
		cable = append(cable, seg)
	}
	return cable, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("input")
	checkErr(err)
	cont, err := ioutil.ReadAll(f)
	checkErr(err)
	lines := strings.Split(string(cont), "\n")
	first, err := assembleCircuit(lines[0])
	checkErr(err)
	second, err := assembleCircuit(lines[1])
	checkErr(err)
	intersections := first.Intersections(second)
	cl := -1
	for _, in := range intersections {
		new := in.Coord.Distance(coord{0, 0})
		if new != 0 && (new < cl || cl == -1) {
			cl = new
		}
	}
	fmt.Printf("Part 1: %d\n", cl)

	cl = -1
	for _, in := range intersections {
		if in.Steps != 0 && (cl == -1 || in.Steps < cl) {
			cl = in.Steps
		}
	}
	fmt.Printf("Part 2: %d\n", cl)
}
