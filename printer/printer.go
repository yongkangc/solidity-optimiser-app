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

// DFS traversal of the AST
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

func (p *Printer) traverseStructMembers(structMembers []ast.Node[ast.NodeType], depth int) {
	for _, member := range structMembers {
		if p.visited[member] {
			continue
		}
		p.visited[member] = true
		p.visitNode(member, depth)
		childNodes := member.GetNodes()
		if len(childNodes) > 0 {
			p.traverseStructMembers(childNodes, depth+1)
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
	// Return because this would give us duplicate printing. however we could use this as well. 
	return
}

func (p *Printer) VisitPragma(pragma *ast.Pragma, depth int) {
	p.output += fmt.Sprintf("%s%s\n", strings.Repeat("  ", depth), pragma.GetText())
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
	p.output += fmt.Sprintf("%sContract %s", indent, contract.Name)
	if contract.BaseContracts != nil {
		// p.output += " is "
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
	indent := strings.Repeat("    ", depth)
	p.output += fmt.Sprintf("%sstruct %s {\n", indent, structDef.Name)

	// traverse the fields of the struct
	p.traverseNodes(structDef.Members, depth+1)
	p.output += fmt.Sprintf("%s}\n", indent)
}

func (p *Printer) VisitVariableDeclaration(variable *ast.Parameter, depth int) {
	indent := strings.Repeat("    ", depth)

	typeName := variable.TypeName.Name
	// visit the type of the variable
	p.visited[variable.TypeName] = true

	p.output += fmt.Sprintf("%s%s %s; \n", indent, typeName, variable.Name)
}

func (p *Printer) VisitTypeName(typeName *ast.TypeName, depth int) {
	p.output += fmt.Sprintf("%sTypeName: %s\n", strings.Repeat("  ", depth), typeName.Name)
}

