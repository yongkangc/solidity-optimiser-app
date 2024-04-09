package printer

import (
	"context"
	"os"

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
