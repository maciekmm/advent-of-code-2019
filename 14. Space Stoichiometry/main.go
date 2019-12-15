package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type recipe struct {
	result      ingredient
	ingredients []ingredient
}

type ingredient struct {
	id       string
	quantity int
}

func parseIngredient(ingr string) (*ingredient, error) {
	raw := strings.Split(ingr, " ")
	if len(raw) != 2 {
		return nil, errors.New("invalid number of igredients, expected {Q} {ID}, got " + ingr)
	}
	q, err := strconv.Atoi(raw[0])
	if err != nil {
		return nil, err
	}
	return &ingredient{
		id:       raw[1],
		quantity: q,
	}, nil
}

func parseRecipe(rawRecipe string) (*recipe, error) {
	inOut := strings.Split(rawRecipe, " => ")
	if len(inOut) != 2 {
		return nil, errors.New("invalid number of parameters, expected {Q} {ID}, got " + rawRecipe)
	}
	rawIngredients := strings.Split(inOut[0], ", ")
	ingredients := make([]ingredient, len(rawIngredients))
	for i, rawIngredient := range rawIngredients {
		ingr, err := parseIngredient(rawIngredient)
		if err != nil {
			return nil, err
		}
		ingredients[i] = *ingr
	}
	output, err := parseIngredient(inOut[1])
	if err != nil {
		return nil, err
	}
	return &recipe{
		result:      *output,
		ingredients: ingredients,
	}, nil
}

func oresRequired(of string, quantity int, recipes map[string]recipe, leftovers map[string]int) int {
	if val, ok := leftovers[of]; ok {
		if val > quantity {
			leftovers[of] = val - quantity
			return 0
		}
		delete(leftovers, of)
		quantity -= val
	}
	if of == "ORE" {
		return quantity
	}
	recip := recipes[of]
	batchesToProduce := int(math.Ceil(float64(quantity) / float64(recip.result.quantity)))
	lefts := recip.result.quantity*batchesToProduce - quantity
	ores := 0
	for _, ingr := range recip.ingredients {
		ores += oresRequired(ingr.id, ingr.quantity*batchesToProduce, recipes, leftovers)
	}
	if lefts != 0 {
		leftovers[of] += lefts
	}
	return ores
}

func main() {
	file, err := ioutil.ReadFile("input")
	checkErr(err)
	rawRecipes := strings.Split(string(file), "\n")
	recipes := make(map[string]recipe)
	for _, rawRecipe := range rawRecipes {
		parsedRecipe, err := parseRecipe(rawRecipe)
		checkErr(err)
		recipes[parsedRecipe.result.id] = *parsedRecipe
	}
	leftovers := make(map[string]int)
	ores := oresRequired("FUEL", 1, recipes, leftovers)
	fmt.Println("Part 1: ", ores)
	supplied := 1000000000000
	ableToProduce := supplied / ores
	step := 10_000

	for step != 0 {
		leftovers := map[string]int{"ORE": supplied}
		ores := oresRequired("FUEL", ableToProduce, recipes, leftovers)
		if ores == 0 {
			ableToProduce += step
		} else {
			step /= 2
			ableToProduce -= step
		}
	}
	fmt.Println("Part 2: ", ableToProduce-1)
}
