package strategy

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/sunmao/strategy/filter"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"

	"github.com/uopensail/recgo-engine/strategy/freqfilter"
	"github.com/uopensail/recgo-engine/strategy/insert"
	"github.com/uopensail/recgo-engine/strategy/rank"
	"github.com/uopensail/recgo-engine/strategy/scatter"

	"github.com/uopensail/recgo-engine/strategy/recalls"
	"github.com/uopensail/recgo-engine/strategy/weighted"
	"github.com/uopensail/recgo-engine/userctx"
)

func init() {
	RegisterPlugin(DefalutStrategy, NewDefaultStrategy)
}

type DefaultStrategyEntity struct {
	cfg table.StrategyEntityMeta

	FilterGroupEntity *freqfilter.FilterGroupEntity

	RecallGroupEntity *recalls.RecallGroupEntity
	WeightedEntity    weighted.IStrategyEntity
	RankEntity        rank.IStrategyEntity
	ScatterEntity     scatter.IStrategyEntity
	InsertGroupEntity *insert.InsertGroupEntity
}

func NewDefaultStrategy(cfg table.StrategyEntityMeta, env config.EnvConfig) IStrategyEntity {
	return &DefaultStrategyEntity{
		cfg: cfg,
	}
}
func (entity *DefaultStrategyEntity) Meta() *table.StrategyEntityMeta {
	return &entity.cfg
}

type Filters []filter.IFliter

func (fiters Filters) Check(id int) bool {
	for _, filter := range fiters {
		if filter.Check(id) == false {
			return false
		}
	}
	return true
}
func (entity *DefaultStrategyEntity) Do(uCtx *userctx.UserContext) (model.StageResult, error) {

	//生成过滤器

	iFilters := make(Filters, 0, 2)
	iFilters = append(iFilters, &uCtx.UserFilter)
	fg, err := entity.FilterGroupEntity.Do(uCtx)
	if err != nil {
		zlog.LOG.Error("FilterGroupEntity.Do", zap.Error(err))
	} else {
		iFilters = append(iFilters, fg)
	}
	//召回
	stageRes, _ := entity.RecallGroupEntity.Do(uCtx, iFilters)
	//排序
	if entity.RankEntity != nil {
		stageRes, _ = entity.RankEntity.Do(uCtx, stageRes)
	}

	//加权
	if entity.WeightedEntity != nil {
		stageRes, _ = entity.WeightedEntity.Do(uCtx, stageRes)
	}

	//打散
	if entity.ScatterEntity != nil {
		stageRes, _ = entity.ScatterEntity.Do(uCtx, stageRes)
	}
	//强插
	if entity.InsertGroupEntity != nil {
		stageRes, _ = entity.InsertGroupEntity.Do(uCtx, stageRes)
	}
	return stageRes, nil
}

func BuildRuntimeDefaultStrategyEntity(from IStrategyEntity, entities *ModelEntities, uCtx *userctx.UserContext) IStrategyEntity {
	cloneEntity := DefaultStrategyEntity{
		cfg: *from.Meta(),
	}

	filterGroupMeta := entities.Model.FilterGroupEntityTableModel.Get(cloneEntity.cfg.FilterGroupEntityID)
	fg := freqfilter.BuildRuntimeEntity(&entities.FilterEntities, uCtx.DBTabelModel,
		uCtx, filterGroupMeta)
	cloneEntity.FilterGroupEntity = fg
	insertGroupMeta := entities.Model.InsertGroupEntityTableModel.Get(cloneEntity.cfg.InsertGroupEntityID)
	var insertRecallIDs []int
	if insertGroupMeta != nil {
		cloneEntity.InsertGroupEntity = insert.BuildRuntimeGroupEntity(&entities.InsertEntities,
			&entities.Model, uCtx, insertGroupMeta)
		insertRecallIDs = cloneEntity.InsertGroupEntity.GetRecallIDs()
	}

	recallGroupMeta := entities.Model.RecallGroupEntityTableModel.Get(cloneEntity.cfg.RecallGroupEntityID)
	cloneEntity.RecallGroupEntity = recalls.BuildRuntimeEntity(&entities.RecallEntities,
		&entities.Model, uCtx, recallGroupMeta, insertRecallIDs)

	rankMeta := entities.Model.RankEntityTableModel.Get(cloneEntity.cfg.RankEntityID)
	if rankMeta != nil {
		cloneEntity.RankEntity = rank.BuildRuntimeEntity(&entities.RankEntities, &entities.Model, uCtx, rankMeta)
	}

	weightedMeta := entities.Model.WeightedEntityTableModel.Get(cloneEntity.cfg.WeightedEntityID)
	if weightedMeta != nil {
		cloneEntity.WeightedEntity = weighted.BuildRuntimeEntity(&entities.WeightedEntities, &entities.Model, uCtx, weightedMeta)
	}

	scatterkMeta := entities.Model.ScatterEntityTableModel.Get(cloneEntity.cfg.ScatterEntityID)
	if scatterkMeta != nil {
		cloneEntity.ScatterEntity = scatter.BuildRuntimeEntity(&entities.ScatterEntities, &entities.Model, uCtx, scatterkMeta)
	}

	return &cloneEntity
}
