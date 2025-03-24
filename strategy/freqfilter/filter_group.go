package freqfilter

import (
	"strconv"
	"sync"

	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/strategy/freqfilter/resource"

	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
)

type Block []byte

func (b Block) Put(v int) {
	b[v>>3] |= (1 << (v & 7))
}

func (b Block) Exists(v int) bool {
	return (b[v>>3] & (1 << (v & 7))) != 0
}

type filterBlock struct {
	iFilter IFilterStrategyEntity
	block   Block
}

func newFilterBlock(block Block, iFilter IFilterStrategyEntity) *filterBlock {
	return &filterBlock{
		iFilter: iFilter,
		block:   block,
	}
}

func (f *filterBlock) Build(userId string, pl *pool.Pool, ress *resource.Resources) {

	ids := f.iFilter.Do(userId, ress)

	for i := 0; i < len(ids); i++ {
		item := pl.GetByKey(ids[i])
		if item != nil && item.ID >= 0 {
			f.block.Put(item.ID)
		}
	}
}

type filterRuntime struct {
	filters []*filterBlock
	block   Block
	pl      *pool.Pool
}

func newFilterRuntime(entities []IFilterStrategyEntity, pl *pool.Pool) *filterRuntime {
	block := make([]byte, (pl.Len()>>3)+1)
	filters := make([]*filterBlock, 0, len(entities))
	for i := 0; i < len(entities); i++ {
		filter := newFilterBlock(block, entities[i])
		filters = append(filters, filter)
	}
	return &filterRuntime{
		filters: filters,
		block:   block,
		pl:      pl,
	}
}

func (fg *filterRuntime) build(userId string, ress *resource.Resources) {
	var wg sync.WaitGroup
	wg.Add(len(fg.filters))
	for i := 0; i < len(fg.filters); i++ {
		go func(idx int) {
			defer wg.Done()
			fg.filters[idx].Build(userId, fg.pl, ress)

		}(i)
	}
	wg.Wait()
}

// True:pass False: filted
func (fg *filterRuntime) Check(id int) bool {
	return (id >= 0 && !fg.block.Exists(id))
}

type IFliter interface {
	Check(id int) bool // True:pass False: filted
}

type FilterGroupEntity struct {
	EntiyMeta *table.FilterGroupEntityMeta
	Entities  []IFilterStrategyEntity
}

func (entity *FilterGroupEntity) Meta() *table.FilterGroupEntityMeta {
	return entity.EntiyMeta
}

func (entity *FilterGroupEntity) Do(uCtx *userctx.UserContext) (IFliter, error) {
	stat := prome.NewStat("Filter.Do")
	defer stat.End()
	filterRuntime := newFilterRuntime(entity.Entities, uCtx.Ress.Pool)
	filterRuntime.build(uCtx.UID(), uCtx.FilterRess)

	return filterRuntime, nil
}

func BuildRuntimeEntity(entities *FilterEntities, dbModel *dbmodel.DBTabelModel,
	uCtx *userctx.UserContext, emeta *table.FilterGroupEntityMeta) *FilterGroupEntity {
	if emeta == nil {
		return nil
	}
	//确认是否命中实验
	caseValue := uCtx.UserAB.EvalFeatureValue(emeta.ABLayerID)
	if len(caseValue) > 0 {
		//查找实验变体
		relateID, err := strconv.Atoi(caseValue)
		//abEntiy := entitys.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.FilterGroupEntityTableModel.Get(relateID)
			if expMeta != nil {
				emeta = expMeta
			}
		}
	}
	ret := FilterGroupEntity{
		EntiyMeta: emeta,
		Entities:  make([]IFilterStrategyEntity, 0, len(emeta.FilterEntities)),
	}
	//实时构建子实体
	for i := 0; i < len(emeta.FilterEntities); i++ {
		entityMeta := dbModel.FilterEntityTableModel.Get(emeta.FilterEntities[i])
		entity := BuildFilterEntity(entities, dbModel, uCtx, entityMeta)
		if entity != nil {
			ret.Entities = append(ret.Entities, entity)
		}

	}
	return &ret

}
