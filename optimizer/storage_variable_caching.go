package optimizer

import (
	"fmt"
	"strings"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
)

// caches storage variables in local variables
func (o *Optimizer) optimizeStorageVariableCaching() {
	// TODO: check if the storage variable is read more than once
	//       if it is, cache it in a local variable
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		// iterate through the contract's functions
		functions := contract.GetFunctions()
		for _, f := range functions {
			modifier := f.GetStateMutability()
			// TODO: we can check if the function modifies the state, but that is another feature
			if modifier != ast_pb.Mutability_VIEW {
				continue
			}
			stateVariables := make(map[int64]*ast.StateVariableDeclaration, 0)
			referencesToStateVariables := make(map[int64][]*ast.PrimaryExpression, 0)
			tree := f.GetAST().GetTree()
			tree.ExecuteCustomTypeVisit(f.GetAST().GetNodes(), ast_pb.NodeType_IDENTIFIER, func(node ast.Node[ast.NodeType]) (bool, error) {
				var exp *ast.PrimaryExpression
				var ok bool
				if exp, ok = node.(*ast.PrimaryExpression); !ok {
					return true, nil
				}
				// referenced declaration of 0 referes to a built in function
				if exp.GetReferencedDeclaration() == 0 {
					return true, nil
				}

				isStateVariable := false
				var decl *ast.StateVariableDeclaration

				if _, ok = stateVariables[exp.GetReferencedDeclaration()]; ok {
					isStateVariable = true
					decl = stateVariables[exp.GetReferencedDeclaration()]
				} else {
					if d, ok := tree.GetById(exp.GetReferencedDeclaration()).(*ast.StateVariableDeclaration); ok {
						if d.IsStateVariable() {
							isStateVariable = true
							decl = d
						}
					}
				}

				if isStateVariable && decl != nil {
					stateVariables[decl.GetId()] = decl
					referencesToStateVariables[decl.GetId()] = append(referencesToStateVariables[decl.GetId()], exp)
				}

				return true, nil
			})

			for sv, ref := range referencesToStateVariables {
				if len(ref) < 2 {
					delete(referencesToStateVariables, sv)
				}
			}

			for _, sv := range stateVariables {
				if _, ok := referencesToStateVariables[sv.GetId()]; !ok {
					continue
				}
				// HACK: doing this screws up the numbering of the nodes
				InsertCachedVariable(f.GetBody().Unit, sv)
				for _, ident := range referencesToStateVariables[sv.GetId()] {
					ident.Name = fmt.Sprintf("cached_%s", sv.GetName())
				}
			}
		}
	}
}

// InsertCachedVariable inserts a cached variable declaration at the beginning of the function body
func InsertCachedVariable(body *ast.BodyNode, sv *ast.StateVariableDeclaration) {
	// create a new variable declaration
	cachedName := fmt.Sprintf("cached_%s", sv.GetName())

	// check if the type of the state variable is a reference type
	hasStorageLocation := false
	referenceTypes := []string{"mapping", "array", "struct"}
	for _, t := range referenceTypes {
		if strings.Contains(sv.GetTypeDescription().GetIdentifier(), t) {
			hasStorageLocation = true
			break
		}
	}
	// if the type is a reference type, we need to store it in memory
	loc := ast_pb.StorageLocation_DEFAULT
	if hasStorageLocation {
		loc = ast_pb.StorageLocation_MEMORY
	}

	cachedVarDeclaration := &ast.VariableDeclaration{
		Declarations: []*ast.Declaration{
			{
				Name:            cachedName,
				TypeName:        sv.GetTypeName(),
				StorageLocation: loc,
			},
		},
		NodeType: ast_pb.NodeType_VARIABLE_DECLARATION,
		InitialValue: &ast.PrimaryExpression{
			NodeType:              ast_pb.NodeType_IDENTIFIER,
			Name:                  sv.GetName(),
			ReferencedDeclaration: sv.GetId(),
			TypeName:              sv.GetTypeName(),
		},
	}
	// put the new variable declaration at the beginning of the function body
	body.Statements = append([]ast.Node[ast.NodeType]{cachedVarDeclaration}, body.Statements...)
}
