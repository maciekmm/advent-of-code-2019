package main

import (
	"fmt"
	"io/ioutil"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadImage(width int, height int) ([][]int, error) {
	file, err := ioutil.ReadFile("input")
	layers := [][]int{}
	if err != nil {
		return nil, err
	}
	for i, val := range file {
		layer := i / (width * height)
		for layer >= len(layers) {
			layers = append(layers, []int{})
		}
		layers[layer] = append(layers[layer], int(val)-'0')
	}
	return layers, nil
}

func layerPixelCount(layer []int) map[int]int {
	stats := make(map[int]int)
	for i := 0; i <= 9; i++ {
		stats[i] = 0
	}
	for _, pixel := range layer {
		stats[pixel]++
	}
	return stats
}

func stackLayers(layers [][]int) []int {
	final := make([]int, len(layers[0]))
	for i := range final {
		final[i] = 2
	}
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		for j, val := range layer {
			if val == 2 {
				continue
			}
			final[j] = val
		}
	}
	return final
}

func printLayer(layer []int, width int, height int) {
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			color := layer[h*width+w]
			if color != 1 {
				fmt.Printf("\u001b[30m█")
			} else if color == 1 {
				fmt.Printf("\u001b[37m█")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	layers, err := loadImage(25, 6)
	checkErr(err)

	var fewestZeros map[int]int
	for _, layer := range layers {
		stats := layerPixelCount(layer)
		if fewestZeros == nil || stats[0] < fewestZeros[0] {
			fewestZeros = stats
		}
	}

	fmt.Printf("Part 1: %d\n", fewestZeros[1]*fewestZeros[2])
	fmt.Printf("Part 2: \n")
	printLayer(stackLayers(layers), 25, 6)
}
