package insert

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
)

type posEntity struct {
	IStrategyEntity
	TargetPos int
}

type entityList []posEntity

func (p entityList) Len() int      { return len(p) }
func (p entityList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p entityList) Less(i, j int) bool {
	if p[i].TargetPos == p[j].TargetPos {
		return p[i].Meta().Priority > p[j].Meta().Priority
	}
	return p[i].TargetPos < p[j].TargetPos
}

type InsertGroupEntity struct {
	EntiyMeta *table.InsertGroupEntityMeta
	Entities  []posEntity
}

func (entity *InsertGroupEntity) Meta() *table.InsertGroupEntityMeta {
	return entity.EntiyMeta
}

func (entity *InsertGroupEntity) GetRecallIDs() []int {
	ret := make([]int, 0, len(entity.Entities))
	for i := 0; i < len(entity.Entities); i++ {
		entityMeta := entity.Entities[i].Meta()
		bFind := false
		for j := 0; j < len(ret); j++ {
			if ret[j] == entityMeta.RecallID {
				bFind = true
				break
			}
		}
		if bFind == false {
			ret = append(ret, entityMeta.RecallID)
		}
	}
	return ret
}

func (entity *InsertGroupEntity) Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error) {
	waitInserts := make([][]int, len(entity.Entities))
	wg := sync.WaitGroup{}
	wg.Add(len(entity.Entities))
	for i := 0; i < len(entity.Entities); i++ {

		go func(index int) {
			defer wg.Done()
			tuples, err := entity.Entities[index].Do(uCtx, in)
			if err == nil {
				waitInserts[index] = tuples
			}
		}(i)
	}
	wg.Wait()
	selected := make(map[int]int, len(entity.Entities))
	selectedArray := make([]int, len(entity.Entities))
	for i := 0; i < len(entity.Entities); i++ {
		selectedArray[i] = -1
		for j := 0; j < len(waitInserts[i]); j++ {
			itemListIndex := waitInserts[i][j]
			if itemListIndex < entity.Entities[i].TargetPos {
				//不要把排在前面的往后拉
				continue
			}
			if itemListIndex == entity.Entities[i].TargetPos {
				//目标位置已经满足条件，不需要强插，
				break
			}
			if _, ok := selected[itemListIndex]; !ok {
				//这个候选物料没有被其他强插的策略选中，
				selected[itemListIndex] = i
				selectedArray[i] = itemListIndex
				break
			}
		}
	}

	fromto := make(map[int]int, len(entity.Entities))
	tofrom := make(map[int]int, len(entity.Entities))
	for i := 0; i < len(entity.Entities); i++ {
		targetPos := entity.Entities[i].TargetPos
		selectedItemListIndex := selectedArray[i]

		if selectedItemListIndex >= 0 {
			//这个强插策略找到合适的了
			if _, ok := tofrom[targetPos]; !ok {
				tofrom[targetPos] = selectedItemListIndex
				fromto[selectedItemListIndex] = targetPos
			}
		}
	}
	outStageList := make(model.RecItemList, len(in.StageList))
	j := 0
	for i := 0; i < len(outStageList); i++ {
		if from, ok := tofrom[i]; ok {
			outStageList[i] = in.StageList[from]
			continue
		}
		for j < len(in.StageList) {
			if _, ok := fromto[j]; !ok {
				break
			}
			j++
		}
		outStageList[i] = in.StageList[j]
	}
	zlog.SLOG.Debug("before:", in.StageList, "after: ", outStageList, "insert map", tofrom)
	in.StageList = outStageList
	return in, nil
}

func BuildRuntimeGroupEntity(entities *InsertEntities, dbModel *dbmodel.DBTabelModel,
	uCtx *userctx.UserContext, emeta *table.InsertGroupEntityMeta) *InsertGroupEntity {
	if emeta == nil {
		return nil
	}
	//确认是否命中实验
	expInfo := uCtx.ABData.GetByLayerID(emeta.ABLayerID)
	if expInfo != nil {
		//查找实验变体
		relateID, err := strconv.Atoi(expInfo.CaseValue)
		//abEntiy := entitys.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.InsertGroupEntityTableModel.Get(relateID)
			if expMeta != nil {
				emeta = expMeta
			}
		}
	}

	ret := InsertGroupEntity{
		EntiyMeta: emeta,
		Entities:  make([]posEntity, 0, len(emeta.InsertEntities)),
	}
	//实时构建子实体

	for i := 0; i < len(emeta.InsertEntities); i++ {
		entityID := emeta.InsertEntities[i]
		entityMeta := dbModel.InsertEntityTableModel.Get(entityID)

		if rand.Intn(10000) < int(entityMeta.Prob)*10000 && entityMeta.Bengin <= entityMeta.End {
			targetPos := rand.Intn(entityMeta.End-entityMeta.Bengin+1) + entityMeta.Bengin
			entity := BuildRuntimeEntity(entities, dbModel, uCtx, entityMeta)
			ret.Entities = append(ret.Entities, posEntity{
				IStrategyEntity: entity,
				TargetPos:       targetPos,
			})
		}

	}
	sort.Sort(entityList(ret.Entities))
	return &ret

}
