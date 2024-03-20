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
	o.optimizeStructPacking()
}

func (o *Optimizer) optimizeStructPacking() {
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		// iterate through the contract's structs
		structs := contract.GetStructs()
		for _, s := range structs {
			for _, member := range s.GetMembers() {
				paramName := member.GetName()
				paramType := member.GetType()
				fmt.Printf("%s %s %d\n", paramType, paramName, sizeOf(paramType))
			}
			// TODO: sort the members by size
			// TODO: re-arrange the members in the original struct
			// TODO: check where the struct is used and update the references
		}
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
