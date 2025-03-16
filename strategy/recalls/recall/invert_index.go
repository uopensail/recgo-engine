package recall

import (
	"math"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"

	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
)

type InvertInexRecall struct {
	table.RecallEntityMeta
	table.InvertInexRecallMeta
	condition *resources.Condition
	*pool.Pool
}

func NewInvertInexRecall(meta table.RecallEntityMeta, ress *resources.Resource, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity {
	//

	recall := InvertInexRecall{
		Pool:             ress.Pool,
		RecallEntityMeta: meta,
	}
	recall.InvertInexRecallMeta = meta.ParseInvertInexRecallMeta()
	recall.condition = resources.BuildCondition(ress, ress.Pool.WholeCollection, recall.Condition)

	return &recall
}
func (r *InvertInexRecall) Meta() *table.RecallEntityMeta {
	return &r.RecallEntityMeta
}
func (r *InvertInexRecall) Do(uCtx *userctx.UserContext, ifilter model.IFliter) ([]int, error) {
	stat := prome.NewStat("Recall." + r.EntityMeta.Name)
	defer stat.End()
	ret := r.do(uCtx, ifilter)
	stat.SetCounter(len(ret))
	return ret, nil
}
func (r *InvertInexRecall) formatKeys(userFeats sample.Features) []string {
	// TODO:
	return nil
}

func zigzagMerge(tmpCollectionList []resources.Collection, eachMaxCol int, ifilter model.IFliter) resources.Collection {
	// 预计算总元素数量
	if len(tmpCollectionList) == 0 {
		return []int{}
	}

	// 预计算总元素数量
	total := 0
	maxCols := 0
	for _, row := range tmpCollectionList {
		total += len(row)
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// 预分配内存
	result := make([]int, 0, total)
	if eachMaxCol < 0 {
		eachMaxCol = math.MaxInt
	}
	rowCnt := make([]int, len(tmpCollectionList))
	for i := 0; i < maxCols; i++ {
		for j, row := range tmpCollectionList {
			if i < len(row) && ifilter.Check(row[i]) && rowCnt[j] < eachMaxCol {
				result = append(result, row[i])
				rowCnt[j]++
			}
		}
	}

	return result
}
func (r *InvertInexRecall) do(uCtx *userctx.UserContext, filter model.IFliter) []int {

	keys := r.formatKeys(uCtx.UFeat)
	if len(keys) == 0 {
		return nil
	}
	invertIndexCollection, ok := uCtx.Ress.InvertIndexRess[r.Resource]
	if !ok {
		return nil
	}
	tmpCollectionList := invertIndexCollection.Get(keys)
	if len(tmpCollectionList) <= 0 {
		return nil
	}
	//zigzip merger
	tmpCollection := zigzagMerge(tmpCollectionList, r.EachMaxCol, filter)
	if r.TopK >= 0 {
		topk := r.TopK
		if len(tmpCollection) > topk {
			return tmpCollection[:topk]
		}
	}

	if r.condition != nil {
		tmpCollection = r.condition.Check(uCtx.UFeat, tmpCollection)
	}
	return tmpCollection
}

func (r *InvertInexRecall) Close() {

}
