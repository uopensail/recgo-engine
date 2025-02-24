package recalls

import (
	"sort"
	"strconv"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/strategy/filter"
	"github.com/uopensail/recgo-engine/strategy/recalls/recall"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type RecallGroupEntity struct {
	EntityMeta *table.RecallGroupEntityMeta
	Entities   []recall.IRecallStrategyEntity
}

func (entity *RecallGroupEntity) Do(uCtx *userctx.UserContext, ifilters map[string]filter.IFliter) (model.StageResult, error) {
	stat := prome.NewStat("match.Do")
	defer stat.End()

	result := make(chan model.RecallResult, len(entity.Entities))

	for i := 0; i < len(entity.Entities); i++ {
		matchEntiy := entity.Entities[i]

		go func() {
			filterName := matchEntiy.Meta().DSLMeta.Filter
			var ifilter filter.IFliter
			if v, ok := ifilters[filterName]; ok {
				ifilter = v
			}
			ret, err := matchEntiy.Do(uCtx, ifilter)
			if err != nil {
				zlog.LOG.Error("match.do", zap.Error(err))
			}
			itemList := make(model.ItemScoreList, 0, len(ret))

			for k := 0; k < len(ret); k++ {
				itemSource := uCtx.Pool.GetById(ret[k])
				if itemSource != nil {
					itemList = append(itemList, model.ItemRefScore{
						ItemFeatures: model.ItemFeatures{
							Source:     itemSource,
							MutFeature: sample.NewMutableFeatures(),
						},
						RecallRecord: model.RecallRecord{
							ID: ret[k],
						},
					})
				}
			}
			result <- model.RecallResult{
				Items: itemList,
				Meta:  matchEntiy.Meta(),
			}

		}()

	}

	//处理超时
	//TODO: config Timeout
	ticker := time.NewTicker(time.Duration(50) * time.Second)
	defer ticker.Stop()
	timeout := false
	count := 0
	recallCnt := 0
	results := make([]model.RecallResult, 0, len(entity.Entities))
	for {
		select {
		case ret := <-result:
			results = append(results, ret)
			count++
			recallCnt += len(ret.Items)
		case <-ticker.C:
			timeout = true
		}
		if timeout || count >= len(entity.Entities) {
			break
		}
	}
	stageRet := entity.mergeAndFilter(uCtx, results, recallCnt, true)
	return stageRet, nil
}

func (entity *RecallGroupEntity) mergeAndFilter(uCtx *userctx.UserContext,
	results []model.RecallResult, cnt int, checkFilter bool) model.StageResult {
	itemList := make(model.RecItemList, 0, cnt)
	stageRes := model.StageResult{
		RecallTrace: model.RecallTrace{
			ItemRecallNameMap: make(map[int]*model.ItemRecallTrace, cnt),
			RecallResults:     make([]model.RecallResult, len(results)),
		},
	}
	recallMaxWeight := make(map[int]float32, cnt)

	duped := make(map[int]bool, cnt)
	for i := 0; i < len(results); i++ {
		zlog.SLOG.Debug("recall ", uCtx.UID(), "merge results ", results[i])
		stageRes.RecallResults[i] = results[i]
		rConf := results[i].Meta
		rWeight := entity.EntityMeta.EntityWeights[rConf.ID]
		for j := 0; j < len(results[i].Items); j++ {
			item := results[i].Items[j]
			if item.Source == nil {
				continue
			}

			if _, ok2 := stageRes.ItemRecallNameMap[item.Source.ID]; ok2 == false {
				stageRes.ItemRecallNameMap[item.Source.ID] = &model.ItemRecallTrace{
					RecallIndex: make([]int, 0, len(results)),
				}
			}
			//单路召回去重
			recallIndexs := stageRes.ItemRecallNameMap[item.Source.ID].RecallIndex
			if len(recallIndexs) > 0 && recallIndexs[len(recallIndexs)-1] == i {
				//上一个recallIndex == 当前i说明 单路召回有重复的
				continue
			}
			if _, ok3 := recallMaxWeight[item.ID]; ok3 == false {
				recallMaxWeight[item.ID] = -1e+37
			}
			stageRes.ItemRecallNameMap[item.Source.ID].RecallIndex =
				append(stageRes.ItemRecallNameMap[item.Source.ID].RecallIndex, i)
			if recallMaxWeight[item.ID] < rWeight {
				recallMaxWeight[item.ID] = rWeight
			}
			if _, ok4 := duped[item.ID]; ok4 {
				continue
			}
			duped[item.ID] = true

			item.Score = rWeight
			itemList = append(itemList, item)
		}
	}

	//赋值最大的权重
	for i := 0; i < len(itemList); i++ {
		item := &itemList[i]
		if score, ok := recallMaxWeight[item.ID]; ok {
			item.Score = score
		}

		//增加内置的recallID
		recallNames := stageRes.GetRecallNames(item.ID)
		item.Set("d_c_recall_", &sample.Strings{Value: recallNames})
	}

	sort.Stable(model.ItemScoreList(itemList))
	stageRes.StageList = itemList
	return stageRes
}

func BuildRuntimeEntity(entities *recall.RecallEntities, dbModel *dbmodel.DBTabelModel,
	uCtx *userctx.UserContext, recallGroupMeta *table.RecallGroupEntityMeta, insertRecallIDs []int) *RecallGroupEntity {
	if recallGroupMeta == nil {
		return nil
	}
	//确认是否命中实验
	expInfo := uCtx.ABData.GetByLayerID(recallGroupMeta.ABLayerID)
	if expInfo != nil {
		//查找实验变体
		relateID, err := strconv.Atoi(expInfo.CaseValue)
		//abEntiy := Entities.Model.ABEntityTableModel.Get(int(expInfo.CaseId))
		if err == nil {
			//替换entiyMeta
			expMeta := dbModel.RecallGroupEntityTableModel.Get(relateID)
			if expMeta != nil {
				recallGroupMeta = expMeta
			}
		}
	}
	ret := RecallGroupEntity{
		EntityMeta: recallGroupMeta,
		Entities:   make([]recall.IRecallStrategyEntity, 0, len(recallGroupMeta.RecallEntities)),
	}
	//实时构建子实体
	for i := 0; i < len(recallGroupMeta.RecallEntities); i++ {
		entiyMeta := dbModel.RecallEntityTableModel.Get(recallGroupMeta.RecallEntities[i])
		if entiyMeta != nil {
			recallEntity := recall.BuildRecallEntity(entities, dbModel, uCtx, entiyMeta)
			if recallEntity != nil {
				ret.Entities = append(ret.Entities, recallEntity)
			}
		}
	}
	for i := 0; i < len(insertRecallIDs); i++ {
		entiyMeta := dbModel.RecallEntityTableModel.Get(insertRecallIDs[i])
		if entiyMeta != nil {
			recallEntity := recall.BuildRecallEntity(entities, dbModel, uCtx, entiyMeta)
			if recallEntity != nil {
				ret.Entities = append(ret.Entities, recallEntity)
			}
		}
	}
	return &ret

}
