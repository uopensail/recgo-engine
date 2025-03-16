package recall

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"

	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
)

type ConditionRecall struct {
	table.RecallEntityMeta

	*pool.Pool
	condition *resources.Condition
}

func NewConditionRecall(meta table.RecallEntityMeta, ress *resources.Resource, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity {
	//

	recall := ConditionRecall{
		Pool:             ress.Pool,
		RecallEntityMeta: meta,
	}

	recall.condition = resources.BuildCondition(ress, ress.Pool.WholeCollection, recall.Condition)
	return &recall
}
func (r *ConditionRecall) Meta() *table.RecallEntityMeta {
	return &r.RecallEntityMeta
}
func (r *ConditionRecall) Do(uCtx *userctx.UserContext, ifilter model.IFliter) ([]int, error) {
	stat := prome.NewStat("Recall." + r.EntityMeta.Name)
	defer stat.End()
	ret := r.do(uCtx.UFeat, uCtx.Ress.Pool, ifilter)
	stat.SetCounter(len(ret))
	return ret, nil
}

func (r *ConditionRecall) do(userFeats sample.Features, pl *pool.Pool, filter model.IFliter) []int {

	var tmpCollection []int
	// filter runtime condition
	if r.condition != nil {
		tmpCollection = r.condition.Check(userFeats, tmpCollection)
	}
	k := 0
	for _, v := range tmpCollection {
		if filter.Check(v) {
			tmpCollection[k] = v
			k++
		}
	}
	if k < len(tmpCollection) {
		tmpCollection = tmpCollection[:k]
	}
	return tmpCollection
}

func (r *ConditionRecall) Close() {

}
