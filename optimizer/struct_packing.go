package optimizer

import (
	"fmt"

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
	testBinCompletion()
	// o.optimizeStructPacking()
}

// I found that the problem of struct packing is NP-hard
// It is a 1D variant of the bin packing problem, which is NP-hard
// Instead of using an almost optimal solution,
// we should use exhaustive search to find the optimal solution
// The number of permutations is n! where n is the number of members in the struct
// which is alright in our case since the number of members is small
type Pair struct {
	idx  int
	size int
}

// func (o *Optimizer) binPacking(params []*ir.Parameter) []int {
// 	result := make([]int, len(params))
// 	smallParams := []Pair{}
// 	for i, param := range params {
// 		// move the large parameters to the top of the struct
// 		if sizeOf(param.GetType()) >= 32 {
// 			result = append(result, i)
// 			continue
// 		}
// 		smallParams = append(smallParams, Pair{i, sizeOf(param.GetType())})
// 	}
//
// 	slotSize := 20
// 	res := bestFitDecreasing(smallParams, slotSize)
// 	for i, bin := range res {
// 		fmt.Printf("bin: %d\n", i)
// 		for _, item := range bin.contents {
// 			// print
// 			fmt.Printf("  %d %d\n", item.idx, item.size)
// 		}
// 	}
// 	fmt.Println(len(binCompletion(smallParams, slotSize)))
//
// 	return result
// }

// Bin completion algorithm
// gets the exact solution to the bin packing problem
// ie. the optimal number of bins
// https://cdn.aaai.org/AAAI/2002/AAAI02-110.pdf
// func binCompletion(pairs []Pair, slotSize int) []Slot {
// 	// first we run best fit
// 	approxSolution := bestFitDecreasing(pairs, slotSize)
//
// 	// this gives us an upper bound on the number of bins
// 	// then we check if we have arrived at the optimal solution
// 	// calculate lower bound
// 	// if the lower bound is equal to the upper bound, then we have the exact solution
// 	totalSize := 0
// 	for _, item := range pairs {
// 		totalSize += item.size
// 	}
//
// 	var lowerBound int = int(math.Ceil(float64(totalSize) / float64(slotSize)))
// 	fmt.Printf("lower bound: %d\n", lowerBound)
// 	fmt.Printf("upper bound: %d\n", len(approxSolution))
// 	if len(approxSolution) == lowerBound {
// 		return approxSolution
// 	}
//
// 	// init lower bound bins first
// 	lowerBoundBins := make([]Slot, lowerBound)
// 	for i := 0; i < lowerBound; i++ {
// 		lowerBoundBins[i] = Slot{freeSpace: slotSize, contents: []Pair{}}
// 	}
//
// 	res := packBin(pairs, lowerBoundBins, slotSize, len(approxSolution))
// 	if res == nil {
// 		panic("no solution found")
// 	}
//
// 	fmt.Printf("exact solution for minimum slots: %d\n", len(res))
//
// 	// use branch and bound to find the exact solution
// 	// prune branches when the number of bins >= upper bound
// 	return nil
// }
// func printSlots(slots []Slot) {
// 	for i, bin := range slots {
// 		fmt.Printf("bin: %d , wasted: %d\n", i, bin.freeSpace)
// 		for _, item := range bin.contents {
// 			// print
// 			fmt.Printf("  %d\n", item.size)
// 		}
// 	}
//
// }
//
// func packBin(pairs []Pair, slots []Slot, slotSize int, upperBound int) []Slot {
// 	// prune branches when the number of bins >= upper bound
// 	fmt.Printf("slots: %d\n", len(slots))
// 	if len(slots) > upperBound {
// 		return nil
// 	}
// 	// base case: no more items to pack
// 	if len(pairs) == 0 {
// 		return slots
// 	}
//
// 	pair := pairs[0]
// 	// try to put the item in each slot
//
// 	var bestSolutionSoFar []Slot = nil
//
// 	for i, slot := range slots {
// 		// skip slots that are too small
// 		if pair.size > slot.freeSpace {
// 			continue
// 		}
// 		// if the slot is large enough
// 		// then put the item in the slot
// 		// and recursively call packBin with the remaining items
// 		leftover := pairs[1:]
// 		slotsCopy := deepCopySlotArray(slots)
//
// 		// put the item in the slot
// 		slotsCopy[i].freeSpace -= pair.size
// 		slotsCopy[i].contents = append(slots[i].contents, pair)
// 		// call packBin with the remaining items
// 		res := packBin(leftover, slotsCopy, slotSize, upperBound)
// 		if res == nil {
// 			continue
// 		}
//
// 		if bestSolutionSoFar == nil || len(res) < len(bestSolutionSoFar) {
// 			bestSolutionSoFar = res
// 		}
// 	}
//
// 	if bestSolutionSoFar == nil {
// 		// increase the number of bins
// 		// by creating a new slot
// 		newSlot := Slot{freeSpace: slotSize - pair.size, contents: []Pair{pair}}
// 		return packBin(pairs[1:], append(slots, newSlot), slotSize, upperBound)
// 	}
//
// 	return bestSolutionSoFar
// }

// func bestFitDecreasing(pairs []Pair, slotSize int) []Slot {
// 	// sort the pairs by size
// 	sort.SliceStable(pairs, func(i, j int) bool {
// 		return pairs[i].size > pairs[j].size
// 	})
// 	return bestFit(pairs, slotSize)
// }

func testBinCompletion() {
	slotSize := 10
	// sizes := []int{9, 8, 2, 2, 5, 4}
	// sizes2 := []int{9, 6, 6, 5, 4, 4, 3, 3, 2, 2}
	// binCapacity := 10
	// for i, size := range sizes2 {
	// 	pairs = append(pairs, Pair{i, size})
	// }
	pairs := []Pair{
		{0, 8}, {1, 7}, {2, 6}, {3, 5}, {4, 4}, {5, 3}, {6, 2}, {7, 2}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, {12, 1},
	}
	// binCompletion(pairs, slotSize)

	// 13, [10.2.1] [8.5] [7.6]
	// 12, [10.2] [8.4] [7.5]
}

// https://stackoverflow.com/questions/15660476/bin-packing-exact-np-hard-exponential-algorithm
// func bestFit(pairs []Pair, slotSize int) []Slot {
// 	slots := []Slot{}
// 	slots = append(slots, Slot{freeSpace: slotSize, contents: []Pair{}})
// 	for _, item := range pairs {
// 		minRemainingSpace := slotSize
// 		bestSlot := -1
//
// 		for i, slot := range slots {
// 			// skip slots that are too small
// 			if item.size > slot.freeSpace {
// 				continue
// 			}
// 			remainingSpace := slot.freeSpace - item.size
// 			// keep track of the best fit slot
// 			if remainingSpace < minRemainingSpace {
// 				bestSlot = i
// 				minRemainingSpace = remainingSpace
// 			}
// 		}
//
// 		// if found best fit slot
// 		// then add the item to the slot
// 		if bestSlot != -1 {
// 			slots[bestSlot].freeSpace -= item.size
// 			slots[bestSlot].contents = append(slots[bestSlot].contents, item)
// 			// else create a new slot
// 		} else {
// 			slots = append(slots, Slot{freeSpace: slotSize - item.size, contents: []Pair{item}})
// 		}
// 	}
//
// 	return slots
// }

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
			// o.binPacking(params)

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

// func deepCopySlotArray(src []Slot) []Slot {
// 	// copy slots first
// 	// copy the items in the slots
// 	dst := make([]Slot, len(src))
// 	copy(dst, src)
// 	for i := range src {
// 		dst[i].contents = make([]Pair, len(src[i].contents))
// 		dst[i].freeSpace = src[i].freeSpace
// 		for j := range src[i].contents {
// 			dst[i].contents[j] = src[i].contents[j]
// 		}
// 	}
// 	return dst
// }
