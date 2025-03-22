package insert

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	meta "github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/prome"
)

type InsertEntities struct {
	Entities map[int]IStrategyEntity
}

func NewInsertEntities(newConfs []table.InsertEntityMeta, envCfg config.EnvConfig,
	ress *resources.Resource) *InsertEntities {
	entities := &InsertEntities{
		Entities: make(map[int]IStrategyEntity),
	}
	for _, v := range newConfs {
		s := PluginFactoryCreate(v, envCfg, ress)
		if s != nil {
			entities.Entities[v.ID] = s
		}
	}
	return entities
}

func (entities *InsertEntities) GetStrategy(id int) IStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("match.GetStrategy.%d", id))
	defer stat.End()

	if entity, ok := entities.Entities[id]; ok {
		return entity
	}
	stat.MarkErr()
	return nil

}

func BuildRuntimeEntity(entities *InsertEntities, dbModel *meta.DBTabelModel,
	uCtx *userctx.UserContext, entiyMeta *table.InsertEntityMeta) IStrategyEntity {
	if entiyMeta == nil {
		return nil
	}
	//确认是否命中实验
	expInfo := uCtx.ABData.GetByLayerID(entiyMeta.ABLayerID)
	if expInfo != nil {
		//查找实验变体
		relateID, err := strconv.Atoi(expInfo.CaseValue)
		//abEntiy := Entities.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.InsertEntityTableModel.Get(relateID)
			if expMeta != nil {
				entiyMeta = expMeta
			}
		}
	}
	ret := entities.GetStrategy(entiyMeta.ID)
	return ret

}
