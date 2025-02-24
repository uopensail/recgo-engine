package weighted

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/prome"
)

type WeightedEntities struct {
	Entities map[int]IStrategyEntity
}

func (entities *WeightedEntities) Clone(a *WeightedEntities) {
	entities.Entities = make(map[int]IStrategyEntity)
	if a != nil {
		for k, v := range a.Entities {
			entities.Entities[k] = v
		}
	}
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

func (entities *WeightedEntities) Reload(newConfs []table.WeightedEntityMeta, envCfg config.EnvConfig) {
	oldConfs := make([]table.WeightedEntityMeta, 0, len(entities.Entities))
	for _, v := range entities.Entities {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)

	if len(invalidM)+len(upsertM) <= 0 {
		return
	}

	for k, v := range upsertM {
		s := PluginFactoryCreate(v, envCfg)
		if s != nil {
			entities.Entities[k] = s
		}
	}

	//删除
	for k := range invalidM {
		delete(entities.Entities, k)
	}

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
