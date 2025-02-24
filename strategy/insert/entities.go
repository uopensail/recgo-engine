package insert

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	meta "github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
)

type InsertEntities struct {
	Entities map[int]IStrategyEntity
}

func (entities *InsertEntities) Clone(a *InsertEntities) {
	entities.Entities = make(map[int]IStrategyEntity)
	if a != nil {
		for k, v := range a.Entities {
			entities.Entities[k] = v
		}
	}
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

func (entities *InsertEntities) Reload(newConfs []table.InsertEntityMeta, envCfg config.EnvConfig, pl *pool.Pool) {
	oldConfs := make([]table.InsertEntityMeta, 0, len(entities.Entities))
	for _, v := range entities.Entities {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)

	if len(invalidM)+len(upsertM) <= 0 {
		return
	}

	for k, v := range upsertM {
		s := PluginFactoryCreate(v, envCfg, pl)
		if s != nil {
			if old, ok := entities.Entities[k]; ok {
				old.Close()
			}
			entities.Entities[k] = s
		}
	}

	//删除
	for k := range invalidM {
		if old, ok := entities.Entities[k]; ok {
			old.Close()
		}
		delete(entities.Entities, k)
	}

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
