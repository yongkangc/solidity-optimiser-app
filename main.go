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

	builder, err := printer.GetBuilder(ctx, filepath)
	if err != nil {
		zap.L().Error("Failed to get builder", zap.Error(err))
	}

	zap.L().Info("Parsing and building contract")
	if err := builder.Parse(); err != nil {
		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
	}
	if err := builder.Build(); err != nil {
		zap.L().Error("Failed to build contract", zap.Error(err))
	}

	ast := builder.GetAstBuilder()
	errs := ast.ResolveReferences()
	if len(errs) > 0 {
		zap.L().Error("Failed to resolve references", zap.Errors("resolve errors", errs))
	}

	// rootNode := ast.GetRoot() // for printing
	opt := optimizer.NewOptimizer(builder)
	// opt.PackStructs()
	// opt.OptimizeCallData()
	opt.CacheStorageVariables()
	fmt.Println(builder.GetRoot().Unit.ToSource())
}
