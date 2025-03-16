package recall

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"

	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
)

type RecallEntities struct {
	Entities map[int]IRecallStrategyEntity
}

func NewRecallEntities(newConfs []table.RecallEntityMeta, envCfg config.EnvConfig,
	ress *resources.Resource, dbModel *dbmodel.DBTabelModel) *RecallEntities {
	entities := RecallEntities{
		Entities: make(map[int]IRecallStrategyEntity),
	}
	for k, v := range newConfs {
		s := PluginFactoryCreate(v, ress, dbModel)
		if s != nil {
			entities.Entities[k] = s
		}
	}
	return &entities
}

func (entities *RecallEntities) GetStrategy(id int) IRecallStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("GetStrategy.%d", id))
	defer stat.End()

	if entiy, ok := entities.Entities[id]; ok {
		return entiy
	}
	stat.MarkErr()
	return nil

}

func BuildRecallEntity(entities *RecallEntities, dbModel *dbmodel.DBTabelModel,
	uCtx *userctx.UserContext, entityMeta *table.RecallEntityMeta) IRecallStrategyEntity {
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
			//替换entityMeta
			expMeta := dbModel.RecallEntityTableModel.Get(relateID)
			if expMeta != nil {
				entityMeta = expMeta
			}
		}
	}
	//这里直接从对象池中获取，无需实时创建
	return entities.GetStrategy(entityMeta.ID)
}
