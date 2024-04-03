// Optimizes the contract to use call data instead of memory for external functions.
package optimizer

import (
	"fmt"
	"strings"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
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

func (o *Optimizer) OptimizeCallData() {
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		functions := contract.GetFunctions()
		for _, f := range functions {
			modifier := f.GetStateMutability()
			if modifier == ast_pb.Mutability_PURE || modifier == ast_pb.Mutability_VIEW {
				astParameters := f.GetAST().Parameters.Parameters

				for _, param := range astParameters {
					if !canBeConvertedToCallData(param) {
						continue
					}

					// need to check if parameter is used in the function
					if param.StorageLocation == ast_pb.StorageLocation_MEMORY {
						param.StorageLocation = ast_pb.StorageLocation_CALLDATA
						fmt.Println("Changed storage location to calldata for", param.GetName())
					}
				}
			}
		}
	}
}

// https://docs.soliditylang.org/en/latest/types.html#reference-types
func canBeConvertedToCallData(param *ast.Parameter) bool {
	if param.StorageLocation == ast_pb.StorageLocation_MEMORY {
		return true
	}
	paramType := param.GetTypeName().GetName()
	isSlice := strings.Contains(paramType, "]")
	isMapping := strings.Contains(paramType, "mapping")
	// TODO: we dont handle structs for now
	if isSlice || isMapping {
		return true
	}
	return false
}
