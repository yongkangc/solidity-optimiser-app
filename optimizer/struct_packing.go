package optimizer

import (
	"fmt"
	"optimizer/optimizer/optimizer/binpack"

	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
)

const SLOT_SIZE = 32

type Optimizer struct {
	builder *ir.Builder
	ast     *ast.RootNode
}

func NewOptimizer(builder *ir.Builder) *Optimizer {
	return &Optimizer{
		builder: builder,
		ast:     builder.GetRoot().GetAST(),
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
			members := s.GetAST().GetMembers()
			items := paramsToItems(members)
			optimalSlots := binpack.OptimalBinPacking(items, SLOT_SIZE)

			// re-arrange the members in the original struct
			optimisedParams := make([]ast.Node[ast.NodeType], 0)
			for _, slot := range optimalSlots {
				for _, item := range slot {
					optimisedParams = append(optimisedParams, members[item.Idx])
				}
			}

			// update the struct with the optimised members
			s.GetAST().Members = optimisedParams

			// TODO: check where the struct is used and update the references
		}
		// update the contract with the optimised structs
	}
	// update the ast with the optimised contracts

}

func (o *Optimizer) printParams(params []*ir.Parameter) {
	for _, member := range params {
		paramName := member.GetName()
		paramType := member.GetType()
		fmt.Printf("%s %s %d\n", paramType, paramName, sizeOf(paramType))
	}
}

func sizeOf(paramType string) int {
	if size, ok := sizeMap[paramType]; ok {
		return size
	} else {
		return 32 // type not in map, set to 32 to be safe
	}
}

func paramsToItems(p []*ast.Parameter) []binpack.Item {
	items := make([]binpack.Item, len(p))
	for i, param := range p {
		param.GetType()
		size := sizeOf(param.GetTypeName().GetName())
		item := binpack.Item{
			Idx:  i,
			Size: size,
		}
		items[i] = item
	}
	return items
}
