package binpack

import (
	"sort"
)

type Item struct {
	Idx  int
	Size int
}

type Slot []Item

/*
* OptimalBinPacking takes a list of pairs and a bin capacity and returns a list of slots
* that minimizes the number of slots used to store all pairs.
 */
func OptimalBinPacking(pairs []Item, binCapacity int) []Slot {
	// Sort the pairs in decreasing order of their sizes
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Size > pairs[j].Size
	})

	// Generate all possible solutions
	allSolutions := generateAllSolutions(pairs, binCapacity)

	// Find the solution with the minimum number of slots
	minSolution := findMinSolution(allSolutions)

	return minSolution
}

func generateAllSolutions(pairs []Item, binCapacity int) [][]Slot {
	var solutions [][]Slot
	generateAllSolutionsRecursive(pairs, binCapacity, []Slot{}, &solutions)
	return solutions
}

func generateAllSolutionsRecursive(pairs []Item, binCapacity int, currentSolution []Slot, solutions *[][]Slot) {
	// Base case: If all pairs are processed
	if len(pairs) == 0 {
		*solutions = append(*solutions, append([]Slot{}, currentSolution...))
		return
	}

	// Try to include the current pair in each existing slot
	for i := range currentSolution {
		if sumSizes(currentSolution[i])+pairs[0].Size <= binCapacity {
			newSolution := append([]Slot{}, currentSolution...)
			newSolution[i] = append(newSolution[i], pairs[0])
			generateAllSolutionsRecursive(pairs[1:], binCapacity, newSolution, solutions)
		}
	}

	// Create a new slot with the current pair
	newSolution := append([]Slot{}, currentSolution...)
	newSolution = append(newSolution, []Item{pairs[0]})
	generateAllSolutionsRecursive(pairs[1:], binCapacity, newSolution, solutions)
}

func findMinSolution(solutions [][]Slot) []Slot {
	minSolution := solutions[0]
	minLen := len(minSolution)

	for _, solution := range solutions {
		if len(solution) < minLen {
			minSolution = solution
			minLen = len(solution)
		}
	}

	return minSolution
}

func sumSizes(pairs []Item) int {
	sum := 0
	for _, pair := range pairs {
		sum += pair.Size
	}
	return sum
}
