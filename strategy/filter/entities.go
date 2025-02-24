package filter

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	meta "github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/strategy/filter/resource"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/prome"
)

type IFilterStrategyEntity interface {
	Do(userID string, ress *resource.Resources) []string
	Meta() *table.FilterEntityMeta
	Close()
}
type FilterEntities struct {
	Entities map[int]IFilterStrategyEntity
}

func (entities *FilterEntities) Clone(a *FilterEntities) {
	entities.Entities = make(map[int]IFilterStrategyEntity)
	if a != nil {
		for k, v := range a.Entities {
			entities.Entities[k] = v
		}
	}
}

func (entities *FilterEntities) GetStrategy(id int) IFilterStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("GetStrategy.%d", id))
	defer stat.End()

	if entiy, ok := entities.Entities[id]; ok {
		return entiy
	}
	stat.MarkErr()
	return nil
}

func (entities *FilterEntities) Reload(newConfs []table.FilterEntityMeta, envCfg config.EnvConfig) {
	oldConfs := make([]table.FilterEntityMeta, 0, len(entities.Entities))
	for _, v := range entities.Entities {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)

	if len(invalidM)+len(upsertM) <= 0 {
		return
	}

	for k, v := range upsertM {
		s := NewFilterEntity(v)
		if s != nil {
			// close old
			if old, ok := entities.Entities[k]; ok {
				old.Close()
			}
			entities.Entities[k] = s
		}
	}

	//删除
	for k := range invalidM {
		// close old
		if old, ok := entities.Entities[k]; ok {
			old.Close()
		}
		delete(entities.Entities, k)
	}

}

func BuildFilterEntity(entities *FilterEntities, dbModel *meta.DBTabelModel,
	uCtx *userctx.UserContext, entityMeta *table.FilterEntityMeta) IFilterStrategyEntity {
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
			expMeta := dbModel.FilterEntityTableModel.Get(relateID)
			if expMeta != nil {
				entityMeta = expMeta
			}
		}
	}
	//这里直接从对象池中获取，无需实时创建
	return entities.GetStrategy(entityMeta.ID)
}
