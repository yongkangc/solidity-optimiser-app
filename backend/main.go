package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"optimizer/optimizer/logger"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"
	"os/exec"
	"strings"
	"text/template"

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
	r.POST("/estimate", estimateHandler)

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

	opt := optimizer.NewOptimizer(builder)
	originalCode, ok := ast_printer.Print(rootNode.GetSourceUnits()[0])

	contractName, err := printer.GetContractName(originalCode)
	if err != nil {
		zap.L().Error("Failed to get contract name", zap.Error(err))
	}

	// Rename the contract to Unoptimized
	unoptimized := renameContract(originalCode, contractName, "Unoptimized")
	if !ok {
		zap.L().Error("Error while printing Original AST")
	}

	// write unoptimized code to file system
	if err := ioutil.WriteFile("../estimator/src/unoptimized.sol", []byte(unoptimized), 0644); err != nil {
		zap.L().Error("Failed to write unoptimized code to file system", zap.Error(err))
	}

	// Optimize the contract
	optimizeContract(opt, input.Options)

	// Print optimised AST
	optimizedCode, ok := ast_printer.Print(rootNode.GetSourceUnits()[0])
	if !ok {
		// error
		zap.L().Error("Error while printing Optimised AST")
	}

	// Rename contract to Optimized
	optimized := renameContract(optimizedCode, contractName, "Optimized")

	// write unoptimized code to file system
	if err := ioutil.WriteFile("../estimator/src/optimized.sol", []byte(optimized), 0644); err != nil {
		zap.L().Error("Failed to write optimized code to file system", zap.Error(err))
	}
	c.JSON(http.StatusOK, gin.H{"optimizedCode": optimized, "unoptimizedCode": unoptimized})
}

func estimateHandler(c *gin.Context) {
	zap.L().Info("Estimate handler")

	var input struct {
		TestCode string `json:"testCode"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tryTestFile(input.TestCode)

	// run make on the test and send the output to the client

	cmd := exec.Command("make", "run")
	cmd.Dir = "../estimator"
	var out bytes.Buffer
	cmd.Stdout = &out
	var errout bytes.Buffer
	cmd.Stderr = &errout

	if err := cmd.Run(); err != nil {
		zap.L().Error("Failed to run make", zap.Error(err), zap.String("stderr", errout.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"output": out.String()})
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

func tryTestFile(test string) {
	type testStruct struct {
		Test         string
		ContractName string
		FileName     string
	}

	tmplFile := "test.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	optimizedSb := strings.Builder{}
	if err := tmpl.Execute(&optimizedSb, testStruct{
		Test:         test,
		ContractName: "Optimized",
		FileName:     "optimized.sol",
	}); err != nil {
		zap.L().Error("Failed to execute template", zap.Error(err))
	}
	ioutil.WriteFile("../estimator/test/optimized.t.sol", []byte(optimizedSb.String()), 0644)

	unoptimizedSb := strings.Builder{}
	if err := tmpl.Execute(&unoptimizedSb, testStruct{
		Test:         test,
		ContractName: "Unoptimized",
		FileName:     "unoptimized.sol",
	}); err != nil {
		zap.L().Error("Failed to execute template", zap.Error(err))
	}
	fmt.Println(unoptimizedSb.String())
	ioutil.WriteFile("../estimator/test/unoptimized.t.sol", []byte(unoptimizedSb.String()), 0644)
}

func renameContract(contract string, oldName string, newName string) string {
	return strings.ReplaceAll(contract, oldName, newName)
}
