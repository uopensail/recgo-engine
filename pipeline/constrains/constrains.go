package constrains

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
)

type IConstrains interface {
	Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection
}

type Constains struct {
	scatter *Scatter
	weights []*WeightAdjust
	inserts []*FixedPositionInsert
}

func NewConstains(confs []model.IConstrain) *Constains {
	scatters := make([]*model.ScatterBasedConstrainConfigure, 0, 8)
	weights := make([]*WeightAdjust, 0, 8)
	inserts := make([]*FixedPositionInsert, 0, 8)
	for _, conf := range confs {
		if conf.GetType() == model.ConstraintTypeScatter {
			scatters = append(scatters, conf.(*model.ScatterBasedConstrainConfigure))
		} else if conf.GetType() == model.ConstraintTypeFixedPosition {
			inserts = append(inserts, NewFixedPositionInsert(conf.(*model.FixedPositionInsertedConstrainConfigure)))
		} else if conf.GetType() == model.ConstraintTypeWeightAdjusted {
			weights = append(weights, NewWeightAdjust(conf.(*model.WeightAdjustedConstrainConfigure)))
		}
	}

	return &Constains{
		scatter: NewScatter(scatters),
		weights: weights,
		inserts: inserts,
	}
}

func (c *Constains) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	tmp := collection
	for _, w := range c.weights {
		tmp = w.Do(uCtx, tmp)
	}
	tmp = c.scatter.Do(uCtx, tmp)
	for _, w := range c.inserts {
		tmp = w.Do(uCtx, tmp)
	}
	return tmp
}
