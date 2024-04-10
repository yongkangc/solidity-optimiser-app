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
	// filepath := cwd + "/examples/unoptimized_contracts/struct_packing.sol"
	// filepath := cwd + "/examples/unoptimized_contracts/calldata.sol"
	filepath := cwd + "/examples/unoptimized_contracts/storage_variable_caching.sol"

	detector, _ := printer.GetDetector(ctx, filepath)

	zap.L().Info("Parsing and building contract")
	if err := detector.Parse(); err != nil {
		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
	}
	if err := detector.Build(); err != nil {
		zap.L().Error("Failed to build contract", zap.Error(err))
	}

	ast := detector.GetAST()
	errs := ast.ResolveReferences()
	if len(errs) > 0 {
		zap.L().Error("Failed to resolve references", zap.Errors("resolve errors", errs))
	}

	// Create a new Printer
	// printer_new := printer.New()

	rootNode := ast.GetRoot()
	zap.L().Info("=============================================")
	fmt.Println(rootNode.ToSource())
	// Print the AST
	// printer_new.Print(rootNode)
	// fmt.Println(printer_new.Output())

	// optimize the contract (still in progress)
	opt := optimizer.NewOptimizer(detector.GetIR())
	// opt.OptimizeCallData()
	opt.CacheStorageVariables()

	zap.L().Info("=============================================")
	fmt.Println(rootNode.ToSource())
	// opt.PackStructs()
	// opt.CacheStorageVariables()
	//
	// // Print the optimized contract
	// printer_opt := printer.New()
	// printer_opt.Print(rootNode)
	// fmt.Println(printer_opt.Output())
}
