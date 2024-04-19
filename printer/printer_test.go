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
func TestEmptyContractPrinter(t *testing.T) {
	emptyContractAstRootNode := setupEmptyContractAST(t) // You need to implement setupEmptyContractAST

	t.Run("Empty Contract Solidity File", func(t *testing.T) {
		printer := printer.New()
		printer.Print(emptyContractAstRootNode)

		expectedOutput := `pragma solidity ^0.8.0;ContractEmpty{}`
		assert.Equal(t, TrimWhitespace(expectedOutput), TrimWhitespace(printer.Output()))
	})
}

func TestMultipleContractsPrinter(t *testing.T) {
	multiContractAstRootNode := setupMultiContractAST(t) // You need to implement setupMultiContractAST

	t.Run("Multiple Contracts Solidity File", func(t *testing.T) {
		printer := printer.New()
		printer.Print(multiContractAstRootNode)

		expectedOutput := `pragmasolidity^0.8.0;ContractBase{uint256x;}pragmasolidity^0.8.0;ContractDerivedBase{uint256y;}`
		assert.Equal(t, TrimWhitespace(expectedOutput), TrimWhitespace(printer.Output()))
	})
}

func TrimWhitespace(input string) string {
	// remove all whitespace from the input string
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")
	return input
}

// Utility functions for setup

func setUpAst(t *testing.T, filePath string) *ast.RootNode {
	// Get the current folder path
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	currentFolderPath := filepath.Dir(testFilePath)
	// Create the path to the "struct_packing.sol" file relative to the current folder
	solidityFilePath := filepath.Join(currentFolderPath, filePath)

	t.Log("\nSolidity folder path: ", solidityFilePath)

	// Create a struct AST
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder, err := printer.GetBuilder(ctx, solidityFilePath)
	if err != nil {
		t.Fatalf("Failed to get builder: %v", err)
	}

	if err := builder.Parse(); err != nil {
		t.Fatalf("Failed to parse contract: %v", err)
	}

	if err := builder.Build(); err != nil {
		t.Fatalf("Failed to build contract: %v", err)
	}

	ast := builder.GetAstBuilder()
	return ast.GetRoot()
}

func setupSampleStructAST(t *testing.T) *ast.RootNode {
	return setUpAst(t, "../examples/unoptimized_contracts/struct_packing.sol")
}

func setupEmptyContractAST(t *testing.T) *ast.RootNode {
	return setUpAst(t, "../tests/testdata/Empty.sol")
}

func setupMultiContractAST(t *testing.T) *ast.RootNode {
	return setUpAst(t, "../tests/testdata/MultipleContracts.sol")
}
