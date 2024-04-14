package testing

import (
	"context"
	"fmt"
	"optimizer/optimizer/optimizer"
	"optimizer/optimizer/printer"
	"os"
	"path/filepath"
	"testing"

	"github.com/unpackdev/solgo/printer/ast_printer"
)

const TEST_DIR = "./testdata"

type Options struct {
	filepath    string
	printOutput bool

	calldata        bool
	structpack      bool
	storagevarcache bool
}

func testHelper(options Options) bool {
	if options.filepath == "" {
		fmt.Println("Error: ", "No file path provided")
		return false
	}
	f := filepath.Join(TEST_DIR, options.filepath)
	fmt.Println("Running test on file: ", f)
	// stat to check if file exists
	if _, err := os.Stat(f); os.IsNotExist(err) {
		fmt.Println("Error: ", "File does not exist")
		return false
	}

	// get builder
	ctx := context.Background()
	builder, err := printer.GetBuilder(ctx, f)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	if err := builder.Parse(); err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	if err := builder.Build(); err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	ast := builder.GetAstBuilder()
	errs := ast.ResolveReferences()
	if len(errs) > 0 {
		fmt.Println("Error: ", errs)
		return false
	}
	root := ast.GetRoot()
	unoptimised, ok := ast_printer.Print(root.GetSourceUnits()[0])
	if !ok {
		fmt.Println("Error: ", "Failed to print unoptimised code")
		return false
	}

	opt := optimizer.NewOptimizer(builder)
	// Run the optimiser
	if options.structpack {
		opt.PackStructs()
	}
	if options.calldata {
		opt.OptimizeCallData()
	}
	if options.storagevarcache {
		opt.CacheStorageVariables()
	}

	optimised, ok := ast_printer.Print(root.GetSourceUnits()[0])
	if !ok {
		fmt.Println("Error: ", "Failed to print optimised code")
		return false
	}

	// Check the output
	if options.printOutput {
		fmt.Println("UNOPTIMIZED====================")
		fmt.Println(unoptimised)
		fmt.Println("================================")
		fmt.Println("OPTIMIZED======================")
		fmt.Println(optimised)
		fmt.Println("================================")
	}

	if unoptimised == optimised {
		fmt.Println("Error: ", "Code not optimised")
		return false
	}
	return true
}

func TestOptimiser(t *testing.T) {
	verbose := true
	tests := []Options{
		// {filepath: "struct_packing.sol", printOutput: verbose, calldata: false, structpack: true, storagevarcache: false},
		// {filepath: "storage_variable_caching.sol", printOutput: verbose, calldata: false, structpack: false, storagevarcache: true},
		// {filepath: "calldata.sol", printOutput: verbose, calldata: true, structpack: false, storagevarcache: false},
		{filepath: "OptimizationShowcase.sol", printOutput: verbose, calldata: false, structpack: false, storagevarcache: true},
	}
	for _, test := range tests {
		if testHelper(test) {
			t.Logf("Test passed")
		} else {
			t.Errorf("Test failed")
		}
	}
}
