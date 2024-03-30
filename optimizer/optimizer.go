package optimizer

import "github.com/unpackdev/solgo/ir"

type Optimizer struct {
	builder *ir.Builder
}

func NewOptimizer(builder *ir.Builder) *Optimizer {
	return &Optimizer{
		builder: builder,
	}
}

func (o *Optimizer) Optimize() {
	o.optimizeStructPacking()
}
