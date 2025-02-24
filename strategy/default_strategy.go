package strategy

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/strategy/filter"
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

	filterEntitys     *filter.FilterEntities
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

func (entity *DefaultStrategyEntity) Do(uCtx *userctx.UserContext) (model.StageResult, error) {

	//召回
	iFilters := make(map[string]filter.IFliter)
	for i := 0; i < len(entity.RecallGroupEntity.Entities); i++ {
		//TODO: Build runtime Filter
		recallEntity := entity.RecallGroupEntity.Entities[i]
		recallMeta := recallEntity.Meta()
		for i := 0; i < len(uCtx.DBTabelModel.FilterGroupEntityTableModel.Rows); i++ {
			if uCtx.DBTabelModel.FilterGroupEntityTableModel.Rows[i].Name == recallMeta.DSLMeta.Filter {
				fg := filter.BuildRuntimeEntity(entity.filterEntitys, uCtx.DBTabelModel,
					uCtx, &uCtx.DBTabelModel.FilterGroupEntityTableModel.Rows[i])

				iFilters[recallMeta.DSLMeta.Filter], _ = fg.Do(uCtx)
				break
			}
		}
	}

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
		cfg:           *from.Meta(),
		filterEntitys: &entities.FilterEntities,
	}

	// filterGroupMeta := Entities.Model.FilterGroupEntityTableModel.Get(cloneEntity.cfg.FilterGroupEntityID)

	// cloneEntity.FilterGroupEntity = filter.BuildRuntimeEntity(&Entities.FilterEntities, &Entities.Model, uCtx, filterGroupMeta)
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
