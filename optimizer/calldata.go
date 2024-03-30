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
}

func (o *Optimizer) optimizeCallData() {
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		functions := contract.GetFunctions()
		for _, f := range functions {
			if f.GetVisibility() == ast_pb.Visibility_EXTERNAL {
				astParameters := f.GetAST().Parameters.Parameters

				for _, param := range astParameters {
					paramType := param.GetTypeName().GetName()
					fmt.Println(paramType)
					// Check if the type is valid for storage
					if _, ok := validTypes[paramType]; !ok {
						continue
					}
					if param.StorageLocation == ast_pb.StorageLocation_MEMORY {
						param.StorageLocation = ast_pb.StorageLocation_CALLDATA
					}
					fmt.Println(param.StorageLocation)
				}
				// TODO: Test for Struct and array
			}
		}
	}
}
