package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/printer/ast_printer"
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
		Options      OptimizationConfig `json:"opts"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	builder, err := printer.GetBuilderCode(ctx, input.ContractCode)
	if err != nil {
		zap.L().Error("Failed to get builder", zap.Error(err))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		zap.L().Error("Failed to get detector", zap.Error(err))
		return
	}

	// Parse the contract
	if err := builder.Parse(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		zap.L().Error("Failed to parse contract", zap.Errors("parse errors", err))
		return
	}

	// Build the contract
	if err := builder.Build(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		zap.L().Error("Failed to build contract", zap.Error(err))
		return
	}

	ast := builder.GetAstBuilder()

	// Resolve references
	if err := resolveReferences(ast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		zap.L().Error("Failed to resolve references", zap.Error(err))
		return
	}

	rootNode := ast.GetRoot()

	// Optimize the contract
	opt := optimizer.NewOptimizer(builder)
	unoptimized, ok := ast_printer.Print(rootNode.GetSourceUnits()[0])
	if ok {
		// write unoptimized code to file system
		if err := ioutil.WriteFile("../estimator/src/unoptimized.sol", []byte(unoptimized), 0644); err != nil {
			zap.L().Error("Failed to write unoptimized code to file system", zap.Error(err))
		}
	}

	optimizeContract(opt, input.Options)
	optimisedCode, ok := ast_printer.Print(rootNode.GetSourceUnits()[0])
	if ok {
		// write unoptimized code to file system
		if err := ioutil.WriteFile("../estimator/src/optimized.sol", []byte(optimisedCode), 0644); err != nil {
			zap.L().Error("Failed to write optimized code to file system", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{"optimizedCode": optimisedCode, "unoptimizedCode": unoptimized})

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
