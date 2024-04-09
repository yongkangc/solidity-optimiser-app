// // test file for printer.go
// // iteratively adding in test cases for the printer.go file
package printer_test

import (
	"context"
	"optimizer/optimizer/printer"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpackdev/solgo/ast"
)

func setupSampleStructAST(t *testing.T) *ast.RootNode {
	// Get the current folder path
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	currentFolderPath := filepath.Dir(testFilePath)
	// Create the path to the "struct_packing.sol" file relative to the current folder
	solidityFilePath := filepath.Join(currentFolderPath, "../examples/unoptimized_contracts/struct_packing.sol")

	t.Log("\nSolidity folder path: ", solidityFilePath)

	// Create a struct AST
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	detector, err := printer.GetDetector(ctx, solidityFilePath)
	if err != nil {
		t.Fatalf("Failed to get detector: %v", err)
	}

	if err := detector.Parse(); err != nil {
		t.Fatalf("Failed to parse contract: %v", err)
	}

	if err := detector.Build(); err != nil {
		t.Fatalf("Failed to build contract: %v", err)
	}

	ast := detector.GetAST()
	return ast.GetRoot()
}

func TestStructPrinter(t *testing.T) {
	structAstRootNode := setupSampleStructAST(t)

	t.Run("Sample Struct Solidity File", func(t *testing.T) {
		printer := printer.New()
		printer.Print(structAstRootNode)

		expectedOutput := `  pragma solidity ^0.8.0;
  Contract NotOptimizedStruct {
        struct Employee {
            uint256 id; 
            uint32 salary; 
            uint32 age; 
            bool isActive; 
            address addr; 
            uint16 department; 
        }
  }
`

		assert.Equal(t, TrimWhitespace(expectedOutput), TrimWhitespace(printer.Output()))

	})
}

func TrimWhitespace(input string) string {
	// remove all whitespace from the input string
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")
	return input
}
