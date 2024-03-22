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
	// o.optimizeStructPacking()
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
