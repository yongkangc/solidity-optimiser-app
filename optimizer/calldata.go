// Optimizes the contract to use call data instead of memory for external functions.
package optimizer

import (
	"fmt"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
)

// All possible variables that can have storage and memory in functions
var validTypes = map[string]bool{
	"string":  true,
	"struct":  true,
	"array":   true,
	"mapping": true,
	"enum":    true,
	"uint256": true,
	"uint64":  true,
	"uint32":  true,
	"uint16":  true,
	"uint8":   true,
	"uint":    true,
}

type Visibility int32

const (
	Visibility_V_DEFAULT Visibility = 0
	Visibility_INTERNAL  Visibility = 1
	Visibility_PRIVATE   Visibility = 2
	Visibility_PUBLIC    Visibility = 3
	Visibility_EXTERNAL  Visibility = 4
)

func (o *Optimizer) optimizeCallData() {
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		functions := contract.GetFunctions()
		for _, f := range functions {
			modifier := f.GetStateMutability()
			if modifier == ast_pb.Mutability_PURE || modifier == ast_pb.Mutability_VIEW {
				astParameters := f.GetAST().Parameters.Parameters

				for _, param := range astParameters {
					paramType := param.GetTypeName().GetName()
					if _, ok := validTypes[paramType]; !ok {
						continue
					}
					if param.StorageLocation == ast_pb.StorageLocation_MEMORY {
						param.StorageLocation = ast_pb.StorageLocation_CALLDATA
						fmt.Println("Changed storage location to calldata for", param.GetName())
					}
				}
				// TODO: Test for Struct and array
			}
		}
	}
}
