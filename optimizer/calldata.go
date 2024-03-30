// Optimizes the contract to use call data instead of memory for external functions.
package optimizer

import (
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
				parameters := f.GetParameters()

				for _, param := range parameters {
					paramType := param.Unit.GetTypeName().GetName()
					// Check if the type is valid for storage
					if _, ok := validTypes[paramType]; !ok {
						continue
					}
					if param.Unit.StorageLocation == ast_pb.StorageLocation_MEMORY {
						param.Unit.StorageLocation = ast_pb.StorageLocation_CALLDATA
					}
					// TODO: check if this changes the AST
				}
			}
		}
	}
}
