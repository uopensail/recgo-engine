package weighted

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/prome"
)

type WeightedEntities struct {
	Entities map[int]IStrategyEntity
}

func NewWeightedEntities(newConfs []table.WeightedEntityMeta, envCfg config.EnvConfig) *WeightedEntities {
	entities := &WeightedEntities{
		Entities: make(map[int]IStrategyEntity),
	}
	for k, v := range newConfs {
		s := PluginFactoryCreate(v, envCfg)
		if s != nil {
			entities.Entities[k] = s
		}
	}
	return entities
}
func (entities *WeightedEntities) GetStrategy(id int) IStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("match.GetStrategy.%d", id))
	defer stat.End()

	if entity, ok := entities.Entities[id]; ok {
		return entity
	}
	stat.MarkErr()
	return nil

}

func BuildRuntimeEntity(entities *WeightedEntities, dbModel *dbmodel.DBTabelModel,
	uCtx *userctx.UserContext, entityMeta *table.WeightedEntityMeta) IStrategyEntity {
	if entityMeta == nil {
		return nil
	}
	//确认是否命中实验
	expInfo := uCtx.ABData.GetByLayerID(entityMeta.ABLayerID)
	if expInfo != nil {
		//查找实验变体
		relateID, err := strconv.Atoi(expInfo.CaseValue)
		//abEntiy := Entities.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.WeightedEntityTableModel.Get(relateID)
			if expMeta != nil {
				entityMeta = expMeta
			}
		}
	}
	ret := entities.GetStrategy(entityMeta.ID)
	return ret

}
