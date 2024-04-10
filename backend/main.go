package main

import (
	"context"
	"net/http"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unpackdev/solgo/ast"
	"go.uber.org/zap"
)

type OptimizationConfig struct {
	StructPacking          bool `json:"structPacking"`
	StorageVariableCaching bool `json:"storageVariableCaching"`
	CallData               bool `json:"callData"`

	// Add more optimization flags here
}

func main() {
	logger.Setup()

	r := gin.Default()
	// Enable CORS
	r.Use(cors.Default())

	r.GET("/health", healthHandler)
	r.POST("/optimize", optimizeHandler)

	r.Run(":8080")
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func optimizeHandler(c *gin.Context) {
	zap.L().Info("Optimize handler")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var input struct {
		ContractCode string             `json:"contractCode"`
		Opts         OptimizationConfig `json:"opts"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detector, err := printer.GetDetectorCode(ctx, input.ContractCode)

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

func optimizeContract(opt *optimizer.Optimizer, config OptimizationConfig) {
	// opt.CacheStorageVariables()
	// TODO: Add more optimization functions here
	// opt.PackStructs()
	if config.StructPacking {
		opt.PackStructs()
	}
	if config.StorageVariableCaching {
		opt.CacheStorageVariables()
	}
	if config.CallData {
		opt.OptimizeCallData()
	}
}
