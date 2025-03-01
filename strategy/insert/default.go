package insert

import (
	"errors"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"

	"github.com/uopensail/recgo-engine/userctx"
)

func init() {
	RegisterPlugin("default", NewDefaultEntity)
}

type DefaultEntity struct {
	ref         utils.Reference
	cfg         table.InsertEntityMeta
	runtimeCond *resources.Condition
}

func NewDefaultEntity(cfg table.InsertEntityMeta, env config.EnvConfig, pl *pool.Pool) IStrategyEntity {
	entity := &DefaultEntity{
		cfg: cfg,
	}
	if len(cfg.Condition) > 0 {
		entity.runtimeCond = resources.BuildCondition(pl, pl.WholeCollection, "pool", cfg.Condition)
	}
	entity.ref.CloseHandler = func() {
		if entity.runtimeCond != nil {
			entity.runtimeCond.Release()
		}
	}
	return entity
}

func (entity *DefaultEntity) Meta() *table.InsertEntityMeta {
	return &entity.cfg
}

func (entity *DefaultEntity) Do(uCtx *userctx.UserContext, in model.StageResult) ([]int, error) {
	collection := make([]int, len(in.StageList))
	itemFeatures := make([]*model.ItemFeatures, len(in.StageList))
	for i := 0; i < len(in.StageList); i++ {
		collection[i] = in.StageList[i].ID
		itemFeatures[i] = &in.StageList[i].ItemFeatures
	}
	if entity.runtimeCond != nil && entity.cfg.Limit > 0 {
		resultC := entity.runtimeCond.CheckWithFillRuntime("user", uCtx.UFeat, collection, "pool", func(id, indexInCollection int) sample.Features {
			return itemFeatures[indexInCollection]
		})
		ret := make([]int, 0, entity.cfg.Limit)
		for i := 0; i < len(resultC); i++ {
			id := resultC[i]
			bFind := -1
			for j := 0; j < len(in.StageList); j++ {
				if id == in.StageList[j].ID {
					bFind = j
					break
				}
			}
			if bFind >= 0 {
				ret = append(ret, bFind)
				if len(ret) >= entity.cfg.Limit {
					break
				}
			}

		}
		if len(ret) > 0 {
			return ret, nil
		}
	}
	return nil, errors.New("force insert empty")
}

func (entity *DefaultEntity) Close() {
	if entity.ref.CloseHandler != nil {
		entity.ref.LazyFree(10)
	}
}
