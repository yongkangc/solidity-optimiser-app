package main

import (
	"context"
	"fmt"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Error("Failed to get current working directory", zap.Error(err))
	}

	// join the current working directory with the file path
	filepath := cwd + "/examples/unoptimized_contracts/struct_packing.sol"

	fmt.Println(filepath)

	detector, _ := printer.GetDetector(ctx, filepath)

	zap.L().Info("Parsing and building contract")
	if err := detector.Parse(); err != nil {
		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
	}
	if err := detector.Build(); err != nil {
		zap.L().Error("Failed to build contract", zap.Error(err))
	}

	ast := detector.GetAST()

	// Create a new Printer
	printer_new := printer.New()

	rootNode := ast.GetRoot()
	// Print the AST
	printer_new.Print(rootNode)
	fmt.Println(printer_new.Output())

	// optimize the contract (still in progress)
	zap.L().Info("Optimizing contract")
	opt := optimizer.NewOptimizer(detector.GetIR())
	opt.Optimize()

	// Print the optimized contract
	printer_opt := printer.New()
	printer_opt.Print(rootNode)
	fmt.Println(printer_opt.Output())
}
