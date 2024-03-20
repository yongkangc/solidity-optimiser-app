package main

import (
	"context"
	"fmt"
	"optimizer/optimizer/logger"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/0x19/solc-switch"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/detector"
	"go.uber.org/zap"
)

func main() {
	logger.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	detector, _ := getDetector(ctx, "./examples/unoptimized_contracts/struct_packing.sol")

	ast := getAST(detector)

	// tree_str, _ := ast.GetTree().ToJSON()
	// zap.S().Infof("AST: %v", string(tree_str))

	// sourceUnits := ast.GetCurrentSourceUnits()
	// jsonStruct, _ := json.Marshal(sourceUnits)

	// zap.S().Infof("Source Units: %v", string(jsonStruct))
	// zap.S().Infof("AST: %v", ast.)

	printCode(ast.GetRoot())

	// // Traverse the AST and print the nodes with dfs
	// for _, node := range tree.GetChildren() {
	// 	zap.S().Infof("Node: %v", node.GetChildCount())
	// }

}

// // printSourceCode prints the source code of the contract.
// func printCode(root *ast.RootNode) {
// 	// zap.S().Infof("Source Units: %v", len(root))
// 	nodes := root.GetNodes()
// 	println(len(nodes))
// 	for _, node := range nodes {
// 		zap.S().Infof("Node: %v", node)
// 		for _, childNode := range node.GetNodes() {
// 			zap.S().Infof("Child Node: %v", childNode)
// 			for _, childChildNode := range childNode.GetNodes() {
// 				zap.S().Infof("childChildNode: %v", childChildNode)
// 				for _, childChildChildNode := range childChildNode.GetNodes() {
// 					zap.S().Infof("childChildChildNode: %v", childChildChildNode)
// 				}
// 			}

// 		}
// 	}
// }

func printCode(root *ast.RootNode) {
	nodes := root.GetNodes()
	println(len(nodes))
	traverseNodes(nodes, 0)
}

func traverseNodes(nodes []ast.Node[ast.NodeType], depth int) {
	for _, node := range nodes {
		printNode(node, depth)
		childNodes := node.GetNodes()
		if len(childNodes) > 0 {
			traverseNodes(childNodes, depth+1)
		}
	}
}

func printNode(node ast.Node[ast.NodeType], depth int) string {
	indent := strings.Repeat("  ", depth)
	var nodeStr string

	switch node.GetType() {
	case ast_pb.NodeType_PRAGMA_DIRECTIVE:
		n := node.(*ast.Pragma)
		nodeStr = fmt.Sprintf("%sPragma: %s", indent, n.GetText())
	// case ast.NodeTypeImport:
	// 	n := node.(*ast.Import)
	// 	nodeStr = fmt.Sprintf("%sImport: %s", indent, n.GetName())
	// case ast.NodeTypeModifierDefinition:
	// 	n := node.(*ast.ModifierDefinition)
	// 	nodeStr = fmt.Sprintf("%sModifier: %s", indent, n.Name)
	// case ast.NodeTypeFunction:
	// 	n := node.(*ast.Function)
	// 	nodeStr = fmt.Sprintf("%sFunction: %s", indent, n.Name)
	// case ast.NodeTypeContract:
	// 	n := node.(*ast.Contract)
	// 	nodeStr = fmt.Sprintf("%sContract: %s", indent, n.Name)
	// case ast.NodeTypeStructDefinition:
	// 	n := node.(*ast.StructDefinition)
	// 	nodeStr = fmt.Sprintf("%sStruct: %s", indent, n.Name)
	// case ast.NodeTypeVariableDeclaration:
	// 	n := node.(*ast.VariableDeclaration)
	// 	nodeStr = fmt.Sprintf("%sVariableDeclaration: %s", indent, n.Src)
	// case ast.NodeTypeTypeName:
	// 	n := node.(*ast.TypeName)
	// 	nodeStr = fmt.Sprintf("%sTypeName: %s", indent, n.Name)
	// case ast.NodeTypeStateVariableDeclaration:
	// 	n := node.(*ast.StateVariableDeclaration)
	// 	nodeStr = fmt.Sprintf("%sStateVariableDeclaration: %s", indent, n.Name)
	// default:
	// 	nodeStr = fmt.Sprintf("%sUnknown Node: %v", indent, node.GetType().String())
	// }
	println(nodeStr)
	return nodeStr
}

func getAST(detector *detector.Detector) *ast.ASTBuilder {
	// Parse the contract
	_ = detector.Parse()

	tree := detector.GetAST()
	return tree
}

// getDetector returns a detector instance for the given file path.
func getDetector(ctx context.Context, filePath string) (*detector.Detector, error) {
	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Error("Failed to get current working directory", zap.Error(err))
		return nil, err
	}

	sources := &solgo.Sources{
		SourceUnits: []*solgo.SourceUnit{
			{
				// Ensure the name matches the contract name. This is crucial!
				Name: "NotOptimizedStruct",
				// Ensure the name in the path matches the contract name. This is crucial!
				Path: filepath.Join(cwd, filePath),
			},
		},
		// Ensure the name matches the base contract name. This is crucial!
		EntrySourceUnitName: "NotOptimizedStruct",
		// Path where additional third party such as openzeppelin are
		LocalSourcesPath: "./examples/",
	}

	config, err := solc.NewDefaultConfig()
	if err != nil {
		zap.L().Error("Failed to construct solc config", zap.Error(err))
		return nil, err
	}

	usr, err := user.Current()
	if err != nil {
		zap.L().Error("Failed to get current user", zap.Error(err))
		return nil, err
	}

	// Make sure that {HOME}/.solc/releases exists prior running this example.
	releasesPath := filepath.Join(usr.HomeDir, ".solc", "releases")
	if err = config.SetReleasesPath(releasesPath); err != nil {
		zap.L().Error("Failed to set releases path", zap.Error(err))
		return nil, err
	}

	compiler, err := solc.New(ctx, config)
	if err != nil {
		zap.L().Error("Failed to construct solc compiler", zap.Error(err))
		return nil, err
	}

	return detector.NewDetectorFromSources(ctx, compiler, sources)
}

// Printer prints out the code based on the AST
type Printer struct {
	output string
}

func NewPrinter() *Printer {
	return &Printer{
		output: "",
	}
}

// Print traverses the AST and prints the nodes
func (p *Printer) Print(node ast.RootNode) {
	// for _, child := range node() {
	// 	p.Print(child)
	// }

}
