// Optimizes the contract to use call data instead of memory for external functions.
package optimizer

import (
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

// if the function argument is only read, it can be converted to calldata
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
					param.StorageLocation = ast_pb.StorageLocation_CALLDATA
				}
			}
		}
	}
}

// https://docs.soliditylang.org/en/latest/types.html#reference-types
func canBeConvertedToCallData(param *ast.Parameter) bool {
	if param.StorageLocation != ast_pb.StorageLocation_MEMORY {
		return false
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

// TODO: Check if the parameter is modified in the function
// might have to do a dfs here to check if the parameter is modified
func isParamModified(param *ast.Parameter) bool {
	return true
}
