// Optimizes the contract to use call data instead of memory for external functions.
package optimizer

import (
	"strings"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
	"go.uber.org/zap"
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
// reference: https://ethereum.stackexchange.com/questions/19380/external-vs-public-best-practices
func (o *Optimizer) OptimizeCallData() {
	zap.L().Info("Optimizing call data")
	nodeVisitor := initVisitor()
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		functions := contract.GetFunctions()
		for _, f := range functions {
			candidates := make([]*ast.Parameter, 0)
			astParameters := f.GetAST().Parameters.Parameters
			for _, param := range astParameters {
				if canBeConvertedToCallData(param) {
					candidates = append(candidates, param)
				}
			}

			modifier := f.GetStateMutability()
			if modifier == ast_pb.Mutability_PURE || modifier == ast_pb.Mutability_VIEW {
				for _, param := range candidates {
					param.StorageLocation = ast_pb.StorageLocation_CALLDATA
				}
				continue
			}

			// check if param is being used in the function
			f.GetAST().GetTree().WalkNode(f.GetAST(), nodeVisitor.visitor)
			for _, param := range candidates {
				if _, found := nodeVisitor.modifiedParams[param.GetName()]; !found {
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

type CallDataVisitor struct {
	modifiedParams map[string]bool
	visitor        *ast.NodeVisitor
}

// TODO: Check if the parameter is modified in the function
// might have to do a dfs here to check if the parameter is modified
func initVisitor() *CallDataVisitor {
	callDataVisitor := &CallDataVisitor{
		modifiedParams: make(map[string]bool, 0),
		visitor:        &ast.NodeVisitor{},
	}
	callDataVisitor.visitor.RegisterTypeVisit(ast_pb.NodeType_ASSIGNMENT, func(node ast.Node[ast.NodeType]) (bool, error) {
		assignment, _ := node.(*ast.Assignment)
		if le := assignment.GetLeftExpression(); le != nil {
			assignment.GetTree().ExecuteCustomTypeVisit(assignment.GetNodes(), ast_pb.NodeType_IDENTIFIER, func(node ast.Node[ast.NodeType]) (bool, error) {
				name := node.(*ast.PrimaryExpression).GetName()
				callDataVisitor.modifiedParams[name] = true
				return true, nil
			})
		}
		// check left hand side
		return true, nil
	})
	return callDataVisitor
}

func isParamModified(param *ast.Parameter, f *ast.Function) bool {
	return false
}
