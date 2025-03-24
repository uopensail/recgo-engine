package freqfilter

import (
	"fmt"
	"strconv"

	"github.com/uopensail/recgo-engine/config"
	meta "github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/strategy/freqfilter/resource"

	"github.com/uopensail/recgo-engine/userctx"

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

func NewFilterEntities(newConfs []table.FilterEntityMeta, envCfg config.EnvConfig) *FilterEntities {
	entities := FilterEntities{
		Entities: make(map[int]IFilterStrategyEntity, len(newConfs)),
	}
	for _, v := range newConfs {
		s := NewFilterEntity(v)
		if s != nil {
			entities.Entities[v.ID] = s
		}
	}
	return &entities
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
func BuildFilterEntity(entities *FilterEntities, dbModel *meta.DBTabelModel,
	uCtx *userctx.UserContext, entityMeta *table.FilterEntityMeta) IFilterStrategyEntity {
	if entityMeta == nil {
		return nil
	}
	//确认是否命中实验
	caseValue := uCtx.UserAB.EvalFeatureValue(entityMeta.ABLayerID)
	if len(caseValue) > 0 {
		//查找实验变体
		relateID, err := strconv.Atoi(caseValue)
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
