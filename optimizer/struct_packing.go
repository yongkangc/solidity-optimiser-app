package optimizer

import (
	"fmt"
	"math"
	"sort"

	"github.com/unpackdev/solgo/ir"
)

type Optimizer struct {
	builder *ir.Builder
}

func NewOptimizer(builder *ir.Builder) *Optimizer {
	return &Optimizer{
		builder: builder,
	}
}

func (o *Optimizer) Optimize() {
	o.optimizeStructPacking()
}

// I found that the problem of struct packing is NP-hard
// It is a 1D variant of the bin packing problem, which is NP-hard
// Instead of using an almost optimal solution,
// we should use exhaustive search to find the optimal solution
// The number of permutations is n! where n is the number of members in the struct
// which is alright in our case since the number of members is small
type pair struct {
	idx  int
	size int
}

func (o *Optimizer) binPacking(params []*ir.Parameter) []int {
	result := make([]int, len(params))
	smallParams := []pair{}
	for i, param := range params {
		// move the large parameters to the top of the struct
		if sizeOf(param.GetType()) >= 32 {
			result = append(result, i)
			continue
		}
		smallParams = append(smallParams, pair{i, sizeOf(param.GetType())})
	}

	slotSize := 20
	res := bestFitDecreasing(smallParams, slotSize)
	for i, bin := range res {
		fmt.Printf("bin: %d\n", i)
		for _, item := range bin.contents {
			// print
			fmt.Printf("  %d %d\n", item.idx, item.size)
		}
	}
	fmt.Println(len(binCompletion(smallParams, slotSize)))

	return result
}

type Slot struct {
	freeSpace int
	contents  []pair
}

// Bin completion algorithm
// gets the exact solution to the bin packing problem
// ie. the optimal number of bins
// https://cdn.aaai.org/AAAI/2002/AAAI02-110.pdf
func binCompletion(pairs []pair, slotSize int) []Slot {
	// first run best fit
	slots := bestFitDecreasing(pairs, slotSize)
	// this gives us an upper bound on the number of bins
	// first check if we have arrived at the optimal solution
	// calculate lower bound
	// if the lower bound is equal to the upper bound, then we have the optimal solution
	totalSize := 0
	for _, item := range pairs {
		totalSize += item.size
	}
	var lowerBound int = int(math.Ceil(float64(totalSize) / float64(slotSize)))
	fmt.Printf("lower bound: %d\n", lowerBound)
	fmt.Printf("upper bound: %d\n", len(slots))
	if len(slots) == lowerBound {
		return slots
	}
	// use branch and bound to find the exact solution
	// prune branches when the number of bins >= upper bound

	return nil
}

func bestFitDecreasing(pairs []pair, slotSize int) []Slot {
	// sort the pairs by size
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].size > pairs[j].size
	})
	return bestFit(pairs, slotSize)
}

// https://stackoverflow.com/questions/15660476/bin-packing-exact-np-hard-exponential-algorithm
func bestFit(pairs []pair, slotSize int) []Slot {
	slots := []Slot{}
	slots = append(slots, Slot{freeSpace: slotSize, contents: []pair{}})
	for _, item := range pairs {
		minRemainingSpace := slotSize
		bestSlot := -1

		for i, slot := range slots {
			// skip slots that are too small
			if item.size > slot.freeSpace {
				continue
			}
			remainingSpace := slot.freeSpace - item.size
			// keep track of the best fit slot
			if remainingSpace < minRemainingSpace {
				bestSlot = i
				minRemainingSpace = remainingSpace
			}
		}

		// if found best fit slot
		// then add the item to the slot
		if bestSlot != -1 {
			slots[bestSlot].freeSpace -= item.size
			slots[bestSlot].contents = append(slots[bestSlot].contents, item)
			// else create a new slot
		} else {
			slots = append(slots, Slot{freeSpace: slotSize, contents: []pair{item}})
		}
	}

	return slots
}

func (o *Optimizer) optimizeStructPacking() {
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		// iterate through the contract's structs
		structs := contract.GetStructs()
		for _, s := range structs {
			o.printParams(s.GetMembers())
			// TODO: sort the members by size
			// deepcopy the members into another array
			params := make([]*ir.Parameter, 0)
			params = append(params, s.GetMembers()...)
			// instead of sorting, we should group the members
			o.binPacking(params)

			// TODO: re-arrange the members in the original struct
			// TODO: check where the struct is used and update the references
		}
	}
}

func (o *Optimizer) printParams(params []*ir.Parameter) {
	for _, member := range params {
		paramName := member.GetName()
		paramType := member.GetType()
		fmt.Printf("%s %s %d\n", paramType, paramName, sizeOf(paramType))
	}
}

// TODO: check if the struct is packed
func (o *Optimizer) isPacked(s *ir.Struct) bool {

	return false
}

func sizeOf(paramType string) int {
	if size, ok := sizeMap[paramType]; ok {
		return size
	} else {
		return 9999 // type not in map, set to huge number
	}
}
