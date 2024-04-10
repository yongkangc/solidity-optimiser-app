package main

import (
	"context"
	"net/http"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"

	"github.com/gin-gonic/gin"
	"github.com/unpackdev/solgo/ast"
	"go.uber.org/zap"
)

func main() {
	logger.Setup()

	r := gin.Default()

	r.POST("/optimize", optimizeHandler)

	r.Run(":8080")
}

func optimizeHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var input struct {
		ContractCode string `json:"contractCode"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detector, err := printer.GetDetector(ctx, input.ContractCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		zap.L().Error("Failed to get detector", zap.Error(err))
		return
	}

	// Parse the contract
	if err := detector.Parse(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
		return
	}

	// Build the contract
	if err := detector.Build(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		zap.L().Error("Failed to build contract", zap.Error(err))
		return
	}

	// Resolve references
	ast := detector.GetAST()
	if err := resolveReferences(ast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		zap.L().Error("Failed to resolve references", zap.Error(err))
		return
	}

	rootNode := ast.GetRoot()

	// Optimize the contract
	opt := optimizer.NewOptimizer(detector.GetIR())
	optimizeContract(opt)

	printer_opt := printer.New()
	printer_opt.Print(rootNode)
	c.JSON(http.StatusOK, gin.H{"optimizedCode": printer_opt.Output()})

}

func resolveReferences(ast *ast.ASTBuilder) error {
	errs := ast.ResolveReferences()
	if len(errs) > 0 {
		zap.L().Error("Failed to resolve references", zap.Errors("resolve errors", errs))
		return errs[0]
	}
	return nil
}

func optimizeContract(opt *optimizer.Optimizer) {
	// opt.CacheStorageVariables()
	// TODO: Add more optimization functions here
	opt.PackStructs()
}

/// Previous code
// func optimize() {
// 	logger.Setup()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		zap.L().Error("Failed to get current working directory", zap.Error(err))
// 	}

// 	// join the current working directory with the file path
// 	// filepath := cwd + "/examples/unoptimized_contracts/struct_packing.sol"
// 	// filepath := cwd + "/examples/unoptimized_contracts/calldata.sol"
// 	filepath := cwd + "/examples/unoptimized_contracts/storage_variable_caching.sol"

// 	detector, _ := printer.GetDetector(ctx, filepath)

// 	zap.L().Info("Parsing and building contract")
// 	if err := detector.Parse(); err != nil {
// 		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
// 	}
// 	if err := detector.Build(); err != nil {
// 		zap.L().Error("Failed to build contract", zap.Error(err))
// 	}

// 	ast := detector.GetAST()
// 	errs := ast.ResolveReferences()
// 	if len(errs) > 0 {
// 		zap.L().Error("Failed to resolve references", zap.Errors("resolve errors", errs))
// 	}

// 	// Create a new Printer
// 	// printer_new := printer.New()

// 	rootNode := ast.GetRoot()
// 	zap.L().Info("=============================================")
// 	fmt.Println(rootNode.ToSource())
// 	// Print the AST
// 	// printer_new.Print(rootNode)
// 	// fmt.Println(printer_new.Output())

// 	// optimize the contract (still in progress)
// 	opt := optimizer.NewOptimizer(detector.GetIR())
// 	// opt.OptimizeCallData()
// 	opt.CacheStorageVariables()

// 	zap.L().Info("=============================================")
// 	fmt.Println(rootNode.ToSource())
// 	// opt.PackStructs()
// 	// opt.CacheStorageVariables()
// 	//
// 	// // Print the optimized contract
// 	// printer_opt := printer.New()
// 	// printer_opt.Print(rootNode)
// 	// fmt.Println(printer_opt.Output())
// }
