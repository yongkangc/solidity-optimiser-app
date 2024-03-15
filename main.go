package main

import (
	"context"
	"os/user"
	"path/filepath"
	"time"

	"github.com/0x19/solc-switch"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/detector"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	currentTick := time.Now()
	defer func() {
		zap.S().Infof("Total time taken: %v", time.Since(currentTick))
	}()

	// Logger setup
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	// Define the Solidity source code for the MyToken contract
	sources := &solgo.Sources{
		SourceUnits: []*solgo.SourceUnit{
			{
				// Ensure the name matches the contract name. This is crucial!
				Name: "NotOptimizedStruct",
				// Ensure the name in the path matches the contract name. This is crucial!
				Path: "NotOptimizedStruct.sol",
				Content: `
        // SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract NotOptimizedStruct {
    struct Employee {
        uint256 id;        // 32 bytes
        uint32 salary;     // 4 bytes
        uint32 age;        // 4 bytes
        bool isActive;     // 1 byte
        address addr;      // 20 bytes
        uint16 department; // 2 bytes
    }
}
`,
			},
		},
		// Ensure the name matches the base contract name. This is crucial!
		EntrySourceUnitName: "NotOptimizedStruct",
		// Path where additional third party such as openzeppelin are
		LocalSourcesPath: "./",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := solc.NewDefaultConfig()
	if err != nil {
		zap.L().Error("Failed to construct solc config", zap.Error(err))
		return
	}

	usr, err := user.Current()
	if err != nil {
		zap.L().Error("Failed to get current user", zap.Error(err))
		return
	}

	// Make sure that {HOME}/.solc/releases exists prior running this example.
	releasesPath := filepath.Join(usr.HomeDir, ".solc", "releases")
	if err = config.SetReleasesPath(releasesPath); err != nil {
		zap.L().Error("Failed to set releases path", zap.Error(err))
		return
	}

	compiler, err := solc.New(ctx, config)
	if err != nil {
		zap.L().Error("Failed to construct solc compiler", zap.Error(err))
		return
	}

	detector, err := detector.NewDetectorFromSources(ctx, compiler, sources)
	if err != nil {
		zap.L().Error("Failed to construct parser", zap.Error(err))
		return
	}

	// Parse the contract
	_ = detector.Parse()

	tree, _ := detector.GetAST().ToJSON()
	str := string(tree)
	zap.S().Infof("AST: %v", str)

	// print ast
	// zap.S().Infof("AST: %v", detector.GetTree().GetText())

	// // Traverse the AST and print the nodes with dfs
	// for _, node := range tree.GetChildren() {
	// 	zap.S().Infof("Node: %v", node.GetChildCount())
	// }

}
