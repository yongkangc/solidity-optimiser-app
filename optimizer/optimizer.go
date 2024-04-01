package optimizer

import (
	"github.com/unpackdev/solgo/ir"
	"go.uber.org/zap"
)

type Optimizer struct {
	builder *ir.Builder
}

func NewOptimizer(builder *ir.Builder) *Optimizer {
	return &Optimizer{
		builder: builder,
	}
}

func (o *Optimizer) PackStructs() {
	zap.L().Info("Packing structs")
	o.optimizeStructPacking()
}

func (o *Optimizer) CacheStorageVariables() {
	zap.L().Info("Caching storage variables")
	o.optimizeStorageVariableCaching()
}
