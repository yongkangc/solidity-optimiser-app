package optimizer

import (
	"fmt"

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

			for _, sv := range stateVariables {
				if !isSupportedType(sv.GetTypeName()) {
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
	cachedVarDeclaration := &ast.VariableDeclaration{
		Id: sv.GetNextID(), // HACK: this is probably bad
		Declarations: []*ast.Declaration{
			{
				Name:     cachedName,
				TypeName: sv.GetTypeName(),
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

// isSupportedType checks if the type is supported for caching
func isSupportedType(t *ast.TypeName) bool {
	name := t.GetName()
	if _, ok := sizeMap[name]; ok {
		return true
	}
	return false
}
