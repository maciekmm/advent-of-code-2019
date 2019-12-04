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

func valid(password int, part2 bool) bool {
	str := strconv.Itoa(password)
	if len(str) != 6 {
		return false
	}
	sameAdjacent := false
	for i := 0; i < len(str)-1; i++ {
		if str[i+1] < str[i] {
			return false
		}
		if str[i+1] == str[i] &&
			(!part2 || ((i == 0 || (i > 0 && str[i-1] != str[i])) && (i >= len(str)-2 || (i < len(str)-2 && str[i+1] != str[i+2])))) {
			sameAdjacent = true
		}
	}
	if !sameAdjacent {
		return false
	}
	return true
}

func main() {
	f, err := os.Open("input")
	checkErr(err)
	cont, err := ioutil.ReadAll(f)
	checkErr(err)
	rang := strings.Split(strings.Trim(string(cont), "\n"), "-")
	from, err := strconv.Atoi(rang[0])
	checkErr(err)
	to, err := strconv.Atoi(rang[1])
	checkErr(err)
	validPwds := 0
	validPwdsPart2 := 0
	for i := from; i <= to; i++ {
		if valid(i, false) {
			validPwds++
		}
		if valid(i, true) {
			validPwdsPart2++
		}
	}
	fmt.Printf("Part 1: %d\n", validPwds)
	fmt.Printf("Part 2: %d\n", validPwdsPart2)
}
