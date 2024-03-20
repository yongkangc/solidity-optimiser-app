package main

import (
	"context"
	"fmt"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/printer"
	"os"
	"os/user"
	"path/filepath"

	"github.com/0x19/solc-switch"
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

	// Create a new Printer
	printer_new := printer.New()

	rootNode := ast.GetRoot()
	// Print the AST
	printer_new.Print(rootNode)
	fmt.Println(printer_new.Output())

	// Previous code
	// printer.PrintCode(rootNode)
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
