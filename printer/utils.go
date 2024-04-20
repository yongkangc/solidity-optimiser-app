package printer

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ir"
	"go.uber.org/zap"
)

// getDetector returns a detector instance for the given file path.

func GetBuilder(ctx context.Context, filePath string) (*ir.Builder, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		zap.L().Error("File does not exist", zap.Error(err))
		return nil, err
	}
	sources := &solgo.Sources{
		SourceUnits: []*solgo.SourceUnit{
			{
				Name: "NotOptimized",
				Path: filePath,
			},
		},
	}
	return ir.NewBuilderFromSources(ctx, sources)
}

func GetBuilderCode(ctx context.Context, code string) (*ir.Builder, error) {
	contractName, err := GetContractName(code)
	if err != nil {
		return nil, err
	}
	sources := &solgo.Sources{
		SourceUnits: []*solgo.SourceUnit{
			{
				Name:    contractName,
				Content: code,
			},
		},
		EntrySourceUnitName: contractName,
	}
	return ir.NewBuilderFromSources(ctx, sources)
}

// getDetector returns a detector instance for the given code

// func GetDetectorCode(ctx context.Context, solidityCode string) (*detector.Detector, error) {
//
//		contractName, err := getContractName(solidityCode)
//
//		if err != nil {
//			return nil, err
//		}
//
//		sources := &solgo.Sources{
//			SourceUnits: []*solgo.SourceUnit{
//				{
//					Name:    contractName,
//					Content: solidityCode,
//				},
//			},
//			EntrySourceUnitName: contractName,
//			// Path where additional third party such as openzeppelin are
//		}
//
//		config, err := solc.NewDefaultConfig()
//		if err != nil {
//			zap.L().Error("Failed to construct solc config", zap.Error(err))
//			return nil, err
//		}
//
//		usr, err := user.Current()
//		if err != nil {
//			zap.L().Error("Failed to get current user", zap.Error(err))
//			return nil, err
//		}
//
//		// Make sure that {HOME}/.solc/releases exists prior running this example.
//		releasesPath := filepath.Join(usr.HomeDir, ".solc", "releases")
//		if err = config.SetReleasesPath(releasesPath); err != nil {
//			zap.L().Error("Failed to set releases path", zap.Error(err))
//			return nil, err
//		}
//
//		compiler, err := solc.New(ctx, config)
//		if err != nil {
//			zap.L().Error("Failed to construct solc compiler", zap.Error(err))
//			return nil, err
//		}
//
//		return detector.NewDetectorFromSources(ctx, compiler, sources)
//	}
//
// // Uses Regex to get the contract name from the code
func GetContractName(solidityCode string) (string, error) {

	// TODO: Deal with case of multiple contracts in a single file
	re := regexp.MustCompile(`contract\s+(\w+)`)
	match := re.FindStringSubmatch(solidityCode)
	fmt.Println(match)

	if len(match) > 1 {
		contractName := match[1]
		return contractName, nil
	} else {
		fmt.Println("Contract not found")
		return "", fmt.Errorf("CONTRACT NOT FOUND")
	}
}
