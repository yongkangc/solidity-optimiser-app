package main

import (
	"context"
	"flag"
	"fmt"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"
	"os"
	"path/filepath"

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

	config := GetConfig()

	// join the current working directory with the file path
	filepath := filepath.Join(cwd, config.filepath)

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

	if config.printOutput {
		fmt.Println("UNOPTIMIZED====================")
		fmt.Println(ast.GetRoot().ToSource())
		fmt.Println("================================")
	}
	opt := optimizer.NewOptimizer(builder)
	if config.packStructs {
		opt.PackStructs()
	}
	if config.optimizeCallData {
		opt.OptimizeCallData()
	}
	if config.cacheStorageVariables {
		opt.CacheStorageVariables()
	}

	if config.printOutput {
		fmt.Println("OPTIMIZED======================")
		fmt.Println(ast.GetRoot().ToSource())
		fmt.Println("================================")
	}
}

type Config struct {
	filepath              string
	packStructs           bool
	optimizeCallData      bool
	cacheStorageVariables bool
	printOutput           bool
}

func GetConfig() Config {
	// use the flag library to parse the command line arguments
	var (
		filepath              string
		packStructs           bool
		optimizeCallData      bool
		cacheStorageVariables bool
		printOutput           bool
	)
	flag.StringVar(&filepath, "file", "", "The path to the file to optimize")
	flag.BoolVar(&packStructs, "pack-structs", false, "Pack structs")
	flag.BoolVar(&optimizeCallData, "optimize-call-data", false, "Optimize call data")
	flag.BoolVar(&cacheStorageVariables, "cache-storage-variables", false, "Cache storage variables")
	flag.BoolVar(&printOutput, "print-output", false, "Print the output")
	flag.Parse()

	fmt.Println("Starting with the following configuration:")
	fmt.Println("  filepath:", filepath)
	fmt.Println("  pack-structs:", packStructs)
	fmt.Println("  optimize-call-data:", optimizeCallData)
	fmt.Println("  cache-storage-variables:", cacheStorageVariables)
	fmt.Println("  print-output:", printOutput)

	if filepath == "" {
		zap.L().Fatal("File path is required")
	}
	return Config{
		filepath:              filepath,
		packStructs:           packStructs,
		optimizeCallData:      optimizeCallData,
		cacheStorageVariables: cacheStorageVariables,
		printOutput:           printOutput,
	}
}
