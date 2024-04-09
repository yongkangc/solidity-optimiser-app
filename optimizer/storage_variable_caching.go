package optimizer

import (
	"fmt"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
)

func (o *Optimizer) optimizeStorageVariableCaching() {
	// TODO: check if the storage variable is read more than once
	//       if it is, cache it in a local variable
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		// iterate through the contract's functions
		functions := contract.GetFunctions()
		for _, f := range functions {
			// iterate through the function's statements
			// f.GetAST().GetTree().Walk(&ast.NodeVisitor{
			// 	Visit: func(node ast.Node[ast.NodeType]) bool {
			// 		if pri, ok := node.(*ast.PrimaryExpression); ok {
			// 			if node.GetType() == ast_pb.NodeType_IDENTIFIER {
			// 				fmt.Println("referenced declaration: ", pri.GetReferencedDeclaration())
			//
			// 				// fmt.Println("source: ", o.builder.GetAstBuilder().GetTree().GetById(pri.GetReferencedDeclaration()).ToSource())
			// 			}
			// 		}
			// 		fmt.Println(node.GetSrc().Line, node.GetId(), node.GetType(), node.ToSource())
			// 		return true
			// 	},
			// })
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
				// HACK: doing this screws up the numbering of the nodes
				InsertCachedVariable(f.GetBody().Unit, sv)
				fmt.Println("Caching state variable: ", sv.ToSource())
				for _, ident := range referencesToStateVariables[sv.GetId()] {
					ident.Name = fmt.Sprintf("cached_%s", sv.GetName())
					fmt.Println("Cached reference: ", ident.GetId(), ident.ToSource())
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
	body.Statements = append([]ast.Node[ast.NodeType]{cachedVarDeclaration}, body.Statements...)
}
