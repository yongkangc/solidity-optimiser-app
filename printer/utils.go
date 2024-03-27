package printer

import (
	"context"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/0x19/solc-switch"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/detector"
	"go.uber.org/zap"
)

// getDetector returns a detector instance for the given file path.

func GetDetector(ctx context.Context, filePath string) (*detector.Detector, error) {

	fmt.Println(filepath.Dir(filePath))
	sources := &solgo.Sources{
		SourceUnits: []*solgo.SourceUnit{
			{
				// Ensure the name matches the contract name. This is crucial!
				Name: "NotOptimizedStruct",
				// Ensure the name in the path matches the contract name. This is crucial!
				Path: filePath,
			},
		},
		// Ensure the name matches the base contract name. This is crucial!
		EntrySourceUnitName: "NotOptimizedStruct",
		// Path where additional third party such as openzeppelin are
		LocalSourcesPath: filepath.Dir(filePath),
	}

	config, err := solc.NewDefaultConfig()
	if err != nil {
		zap.L().Error("Failed to construct solc config", zap.Error(err))
		return nil, err
	}

	usr, err := user.Current()
	if err != nil {
		zap.L().Error("Failed to get current user", zap.Error(err))
		return nil, err
	}

	// Make sure that {HOME}/.solc/releases exists prior running this example.
	releasesPath := filepath.Join(usr.HomeDir, ".solc", "releases")
	if err = config.SetReleasesPath(releasesPath); err != nil {
		zap.L().Error("Failed to set releases path", zap.Error(err))
		return nil, err
	}

	compiler, err := solc.New(ctx, config)
	if err != nil {
		zap.L().Error("Failed to construct solc compiler", zap.Error(err))
		return nil, err
	}

	return detector.NewDetectorFromSources(ctx, compiler, sources)
}
