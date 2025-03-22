package strategy

import (
	"strconv"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/recgo-engine/strategy/freqfilter"
	fresource "github.com/uopensail/recgo-engine/strategy/freqfilter/resource"
	"github.com/uopensail/recgo-engine/strategy/insert"
	"github.com/uopensail/recgo-engine/strategy/rank"
	"github.com/uopensail/recgo-engine/strategy/recalls/recall"

	"github.com/uopensail/recgo-engine/strategy/scatter"
	"github.com/uopensail/recgo-engine/strategy/weighted"
	"github.com/uopensail/recgo-engine/userctx"
)

const (
	DefalutStrategy = "default"
)

type ModelEntities struct {
	Model dbmodel.DBTabelModel // 配置引用

	FilterResources fresource.Resources
	Ress            resources.Resource

	freqfilter.FilterEntities
	recall.RecallEntities
	scatter.ScatterEntities
	insert.InsertEntities
	rank.RankEntities
	weighted.WeightedEntities
	StrategyEntities
}

type IStrategyEntity interface {
	Do(uCtx *userctx.UserContext) (model.StageResult, error)
	Meta() *table.StrategyEntityMeta
}

func BuildRuntimeEntity(entities *ModelEntities, uCtx *userctx.UserContext, entiyMeta *table.StrategyEntityMeta) IStrategyEntity {

	//确认是否命中实验
	expInfo := uCtx.ABData.GetByLayerID(entiyMeta.ABLayerID)
	if expInfo != nil {
		//查找实验变体
		relateID, err := strconv.Atoi(expInfo.CaseValue)
		//abEntiy := Entities.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := entities.Model.StrategyEntityTableModel.Get(relateID)
			if expMeta != nil {
				entiyMeta = expMeta
			}
		}
	}
	cacheEntity := entities.StrategyEntities.GetStrategy(entiyMeta.Name)

	return BuildRuntimeDefaultStrategyEntity(cacheEntity, entities, uCtx)

}
