package recall

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
)

type RecallEntities struct {
	Entities map[int]IRecallStrategyEntity
}

func (entities *RecallEntities) Clone(a *RecallEntities) {
	entities.Entities = make(map[int]IRecallStrategyEntity)
	if a != nil {
		for k, v := range a.Entities {
			entities.Entities[k] = v
		}
	}
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

func (entities *RecallEntities) Reload(newConfs []table.RecallEntityMeta, envCfg config.EnvConfig,
	pl *pool.Pool, poolUpdate bool, dbModel *dbmodel.DBTabelModel) {
	oldConfs := make([]table.RecallEntityMeta, 0, len(entities.Entities))
	for _, v := range entities.Entities {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)
	if poolUpdate {
		//update all
		for _, v := range newConfs {
			upsertM[v.ID] = v
		}
	}
	if len(invalidM)+len(upsertM) <= 0 {
		return
	}

	for k, v := range upsertM {
		s := PluginFactoryCreate(v, pl, dbModel)
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
