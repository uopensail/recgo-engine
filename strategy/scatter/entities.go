package scatter

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	meta "github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/prome"
)

type ScatterEntities struct {
	Entities map[int]IStrategyEntity
}

func NewScatterEntities(newConfs []table.ScatterEntityMeta, envCfg config.EnvConfig) *ScatterEntities {
	entities := &ScatterEntities{
		Entities: make(map[int]IStrategyEntity),
	}
	for _, v := range newConfs {
		s := PluginFactoryCreate(v, envCfg)
		if s != nil {
			entities.Entities[v.ID] = s
		}
	}
	return entities
}

func (entities *ScatterEntities) GetStrategy(id int) IStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("match.GetStrategy.%d", id))
	defer stat.End()

	if entity, ok := entities.Entities[id]; ok {
		return entity
	}
	stat.MarkErr()
	return nil

}

func BuildRuntimeEntity(entities *ScatterEntities, dbModel *meta.DBTabelModel,
	uCtx *userctx.UserContext, entityMeta *table.ScatterEntityMeta) IStrategyEntity {
	if entityMeta == nil {
		return nil
	}
	//确认是否命中实验
	caseValue := uCtx.UserAB.AbInfo.EvalFeatureValue(uCtx.Context, entityMeta.ABLayerID)
	if len(caseValue) > 0 {
		//查找实验变体
		relateID, err := strconv.Atoi(caseValue)
		//abEntiy := Entities.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.ScatterEntityTableModel.Get(relateID)
			if expMeta != nil {
				entityMeta = expMeta
			}
		}
	}
	ret := entities.GetStrategy(entityMeta.ID)
	return ret

}
