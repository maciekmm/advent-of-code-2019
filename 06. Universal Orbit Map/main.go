package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const star = "COM"

var satellites map[string][]string = make(map[string][]string)
var parents map[string][]string = make(map[string][]string)

func orbitters(star string, length int) int {
	orb := 0
	if sats, ok := satellites[star]; ok {
		for _, sat := range sats {
			orb += 1 + length + orbitters(sat, length+1)
		}
	}
	return orb
}

func fewestTransfers(from string, to string) int {
	cnode := from
	visited := make(map[string]interface{})
	length := make(map[string]int)
	length[from] = 0
	queue := append(parents[cnode], satellites[cnode]...)
	for _, b := range queue {
		length[b] = 0
	}
	for len(queue) > 0 {
		body := queue[0]
		queue = queue[1:] // dequeue
		for _, children := range append(parents[body], satellites[body]...) {
			if children == to {
				return length[body]
			}
			if _, visi := visited[children]; !visi {
				queue = append(queue, children)
				length[children] = length[body] + 1
				visited[children] = struct{}{}
			}
		}

	}
	return -1
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
	associations := strings.Split(string(cont), "\n")
	for _, assoc := range associations {
		orb := strings.Split(assoc, ")")
		if ch, ok := satellites[orb[0]]; ok {
			satellites[orb[0]] = append(ch, orb[1])
		} else {
			satellites[orb[0]] = []string{orb[1]}
		}
		if ch, ok := parents[orb[1]]; ok {
			parents[orb[1]] = append(ch, orb[0])
		} else {
			parents[orb[1]] = []string{orb[0]}
		}
	}
	fmt.Printf("%d\n", orbitters(star, 0))
	fmt.Printf("%d\n", fewestTransfers("YOU", "SAN"))
}
