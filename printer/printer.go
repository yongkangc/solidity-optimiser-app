package printer

import (
	"fmt"
	"strings"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
)

// Printer prints out the code based on the AST
type Printer struct {
	output  string
	visited map[ast.Node[ast.NodeType]]bool
}

func New() *Printer {
	return &Printer{
		output:  "",
		visited: make(map[ast.Node[ast.NodeType]]bool),
	}
}

func (p *Printer) Output() string {
	return p.output
}

// Print traverses the AST and prints the nodes
func (p *Printer) Print(root *ast.RootNode) {
	nodes := root.GetNodes()
	p.traverseNodes(nodes, 0)
}

func (p *Printer) traverseNodes(nodes []ast.Node[ast.NodeType], depth int) {
	for _, node := range nodes {
		if p.visited[node] {
			continue
		}
		p.visited[node] = true
		p.visitNode(node, depth)
		childNodes := node.GetNodes()
		if len(childNodes) > 0 {
			p.traverseNodes(childNodes, depth+1)
		}
	}
}

func (p *Printer) visitNode(node ast.Node[ast.NodeType], depth int) {
	switch node.GetType() {
	case ast_pb.NodeType_SOURCE_UNIT:
		n := node.(*ast.SourceUnit[ast.Node[ast_pb.SourceUnit]])
		p.VisitSourceUnit(n, depth)
	case ast_pb.NodeType_PRAGMA_DIRECTIVE:
		n := node.(*ast.Pragma)
		p.VisitPragma(n, depth)
	case ast_pb.NodeType_IMPORT_DIRECTIVE:
		n := node.(*ast.Import)
		p.VisitImport(n, depth)
	case ast_pb.NodeType_MODIFIER_DEFINITION:
		n := node.(*ast.ModifierDefinition)
		p.VisitModifierDefinition(n, depth)
	case ast_pb.NodeType_FUNCTION_DEFINITION:
		n := node.(*ast.Function)
		p.VisitFunctionDefinition(n, depth)
	case ast_pb.NodeType_CONTRACT_DEFINITION:
		n := node.(*ast.Contract)
		p.VisitContractDefinition(n, depth)
	case ast_pb.NodeType_STRUCT_DEFINITION:
		n := node.(*ast.StructDefinition)
		p.VisitStructDefinition(n, depth)
	case ast_pb.NodeType_VARIABLE_DECLARATION:
		n := node.(*ast.Parameter)
		p.VisitVariableDeclaration(n, depth)
	case ast_pb.NodeType_ELEMENTARY_TYPE_NAME:
		n := node.(*ast.TypeName)
		p.VisitTypeName(n, depth)
	default:
		p.output += fmt.Sprintf("Unknown Node: %v\n", node.GetType().String())
		println("%sUnknown Node: %v", node.GetType().String())
	}
}

func (p *Printer) VisitSourceUnit(sourceUnit *ast.SourceUnit[ast.Node[ast_pb.SourceUnit]], depth int) {
	return
}

func (p *Printer) VisitPragma(pragma *ast.Pragma, depth int) {
	p.output += fmt.Sprintf("%sPragma: %s\n", strings.Repeat("  ", depth), pragma.GetText())
}

func (p *Printer) VisitImport(imp *ast.Import, depth int) {
	p.output += fmt.Sprintf("%sImport: %s\n", strings.Repeat("  ", depth), imp.GetName())
}

func (p *Printer) VisitModifierDefinition(modifierDef *ast.ModifierDefinition, depth int) {
	p.output += fmt.Sprintf("%sModifier: %s\n", strings.Repeat("  ", depth), modifierDef.Name)
}

func (p *Printer) VisitFunctionDefinition(function *ast.Function, depth int) {
	p.output += fmt.Sprintf("%sFunction: %s\n", strings.Repeat("  ", depth), function.Name)
}

func (p *Printer) VisitContractDefinition(contract *ast.Contract, depth int) {
	indent := strings.Repeat("  ", depth)
	p.output += fmt.Sprintf("%sContract: %s", indent, contract.Name)
	if contract.BaseContracts != nil {
		p.output += " is "
		for i, base := range contract.BaseContracts {
			if i > 0 {
				p.output += ", "
			}
			p.output += base.BaseName.Name
		}
	}
	p.output += " {\n"
	p.traverseNodes(contract.Nodes, depth+1)
	p.output += fmt.Sprintf("%s}\n", indent)
}

func (p *Printer) VisitStructDefinition(structDef *ast.StructDefinition, depth int) {
	p.output += fmt.Sprintf("%s%s\n", strings.Repeat("  ", depth), structDef.Name)
}

func (p *Printer) VisitVariableDeclaration(variable *ast.Parameter, depth int) {
	p.output += fmt.Sprintf("%s%s\n", strings.Repeat("  ", depth), variable.Name)
}

func (p *Printer) VisitTypeName(typeName *ast.TypeName, depth int) {
	p.output += fmt.Sprintf("%sTypeName: %s\n", strings.Repeat("  ", depth), typeName.Name)
}

/// Previous implementation of the printer

// func printNode(node ast.Node[ast.NodeType], depth int) string {
// 	indent := strings.Repeat("  ", depth)
// 	var nodeStr string

// 	// the idea is to get the type and type cast
// 	switch node.GetType() {
// 	case ast_pb.NodeType_SOURCE_UNIT:
// 		n := node.(*ast.SourceUnit[ast.Node[ast_pb.SourceUnit]])
// 		nodeStr = fmt.Sprintf("%sSourceUnit: %s", indent, n.Name)
// 	case ast_pb.NodeType_PRAGMA_DIRECTIVE:
// 		n := node.(*ast.Pragma)
// 		nodeStr = fmt.Sprintf("%sPragma: %s", indent, n.GetText())
// 	case ast_pb.NodeType_IMPORT_DIRECTIVE:
// 		n := node.(*ast.Import)
// 		nodeStr = fmt.Sprintf("%sImport: %s", indent, n.GetName())
// 	case ast_pb.NodeType_MODIFIER_DEFINITION:
// 		n := node.(*ast.ModifierDefinition)
// 		nodeStr = fmt.Sprintf("%sModifier: %s", indent, n.Name)
// 	case ast_pb.NodeType_FUNCTION_DEFINITION:
// 		n := node.(*ast.Function)
// 		nodeStr = fmt.Sprintf("%sFunction: %s", indent, n.Name)
// 	case ast_pb.NodeType_CONTRACT_DEFINITION:
// 		n := node.(*ast.Contract)
// 		nodeStr = fmt.Sprintf("%sContract: %s", indent, n.Name)
// 	case ast_pb.NodeType_STRUCT_DEFINITION:
// 		n := node.(*ast.StructDefinition)
// 		nodeStr = fmt.Sprintf("%sStruct: %s", indent, n.Name)
// 	case ast_pb.NodeType_VARIABLE_DECLARATION:
// 		n := node.(*ast.Parameter)
// 		nodeStr = fmt.Sprintf("%sVariableDeclaration: %s", indent, n.Name)
// 	case ast_pb.NodeType_ELEMENTARY_TYPE_NAME:
// 		n := node.(*ast.TypeName)
// 		nodeStr = fmt.Sprintf("%sTypeName: %s", indent, n.Name)
// 	default:
// 		// This would give us the type for us to figure out the type of node that is not handled
// 		nodeStr = fmt.Sprintf("%sUnknown Node: %v", indent, node.GetType().String())
// 	}
// 	println(nodeStr)
// 	return nodeStr

// }

// func PrintCode(root *ast.RootNode) {
// 	nodes := root.GetNodes()
// 	traverseNodes(nodes, 0)
// }

// func traverseNodes(nodes []ast.Node[ast.NodeType], depth int) {
// 	for _, node := range nodes {
// 		printNode(node, depth)
// 		childNodes := node.GetNodes()
// 		if len(childNodes) > 0 {
// 			traverseNodes(childNodes, depth+1)
// 		}
// 	}
// }
