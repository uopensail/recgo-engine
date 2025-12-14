package constrains

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IConstrains defines the interface for constraint operations.
// A constraint takes a user context and a collection, and returns a modified collection.
type IConstrains interface {
	Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection
}

// Constains orchestrates execution of multiple constraints:
// 1. Weight adjustments
// 2. Scatter-based distribution
// 3. Fixed position inserts
type Constains struct {
	scatter *Scatter               // Scatter constraint handler
	weights []*WeightAdjust        // List of weight adjustment constraints
	inserts []*FixedPositionInsert // List of fixed position insert constraints
}

// NewConstains constructs a Constains object from a list of constraint configurations.
// It separates different configurations into their respective handlers.
func NewConstains(confs []model.IConstrain) *Constains {
	pStat := prome.NewStat("NewConstains")
	defer pStat.End()
	scatters := make([]*model.ScatterBasedConstrainConfigure, 0, 8)
	weights := make([]*WeightAdjust, 0, 8)
	inserts := make([]*FixedPositionInsert, 0, 8)

	for _, conf := range confs {
		switch conf.GetType() {
		case model.ConstraintTypeScatter:
			if c, ok := conf.(*model.ScatterBasedConstrainConfigure); ok {
				scatters = append(scatters, c)
			} else {
				zlog.LOG.Error("NewConstains.TypeAssertError", zap.String("expected", "ScatterBasedConstrainConfigure"))
			}
		case model.ConstraintTypeFixedPosition:
			if c, ok := conf.(*model.FixedPositionInsertedConstrainConfigure); ok {
				inserts = append(inserts, NewFixedPositionInsert(c))
			} else {
				zlog.LOG.Error("NewConstains.TypeAssertError", zap.String("expected", "FixedPositionInsertedConstrainConfigure"))
			}
		case model.ConstraintTypeWeightAdjusted:
			if c, ok := conf.(*model.WeightAdjustedConstrainConfigure); ok {
				weights = append(weights, NewWeightAdjust(c))
			} else {
				zlog.LOG.Error("NewConstains.TypeAssertError", zap.String("expected", "WeightAdjustedConstrainConfigure"))
			}
		default:
			zlog.LOG.Warn("NewConstains.UnknownType", zap.String("type", conf.GetType()))
		}
	}

	return &Constains{
		scatter: NewScatter(scatters),
		weights: weights,
		inserts: inserts,
	}
}

// Do executes constraints in a fixed sequence:
// 1. Apply all weight adjustments
// 2. Apply scatter-based distributions
// 3. Apply fixed position insertions
func (c *Constains) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("Constains.Do")
	defer pStat.End()
	tmp := collection

	// Apply weight adjustments first
	for _, w := range c.weights {
		tmp = w.Do(uCtx, tmp)
	}

	// Apply scatter distribution next
	tmp = c.scatter.Do(uCtx, tmp)

	// Apply fixed position inserts last
	for _, insert := range c.inserts {
		tmp = insert.Do(uCtx, tmp)
	}

	return tmp
}
