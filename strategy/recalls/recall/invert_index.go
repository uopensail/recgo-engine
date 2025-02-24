package recall

import (
	"sort"

	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/poolsource"
	"github.com/uopensail/recgo-engine/strategy/filter"
	"github.com/uopensail/recgo-engine/strategy/recalls/resource"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
)

type InvertInexRecall struct {
	ref utils.Reference
	table.RecallEntityMeta

	*pool.Pool
}

func NewInvertInexRecall(meta table.RecallEntityMeta, pl *pool.Pool, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity {
	//
	recall := InvertInexRecall{
		Pool:             pl,
		RecallEntityMeta: meta,
	}
	return &recall
}
func (r *InvertInexRecall) Meta() *table.RecallEntityMeta {
	return &r.RecallEntityMeta
}
func (r *InvertInexRecall) Do(uCtx *userctx.UserContext, ifilter filter.IFliter) ([]int, error) {
	stat := prome.NewStat("Recall." + r.EntityMeta.Name)
	defer stat.End()
	ret := r.do(uCtx.UFeat, uCtx.Pool, uCtx.RecallRess, ifilter)
	stat.SetCounter(len(ret))
	return ret, nil
}

func (r *InvertInexRecall) do(userFeats sample.Features, pl *pool.Pool, ress *resource.Resources, filter filter.IFliter) []int {
	var iResource resource.IResource
	if len(r.DSLMeta.Resource) > 0 {
		iResource = ress.Get(r.DSLMeta.Resource)
	}

	if r.DSLMeta.Resource != "pool" && iResource == nil {
		return nil
	}

	var keys []string
	var tmpCollection []int
	if r.formatEval != nil {
		//gen  keys
		keys = r.formatEval.Do(userFeats)
		if len(keys) == 0 {
			return nil
		}
		vms := make(map[int]float32)
		vvs := iResource.Get(keys, pl)
		//merge
		maxJ := len(vvs[0])
		for i := 0; i < len(vvs); i++ {
			if len(vvs[i]) > maxJ {
				maxJ = len(vvs[i])
			}
		}
		keysCollection := make([]int, 0, len(vvs)*maxJ)
		for j := 0; j < maxJ; j++ {
			for i := 0; i < len(vvs); i++ {
				if j < len(vvs[i]) {
					id := vvs[i][j]
					if filter != nil && !filter.Check(id.First) {
						keysCollection = append(keysCollection, id.First)
						if score, ok := vms[id.First]; !ok || id.Second > score {
							vms[id.First] = id.Second
						}

					}
				}
			}
		}

		if len(keysCollection) == 0 {
			return nil
		}

		tmpCollection = poolsource.Intersection(r.Pool.Len(), keysCollection, r.staticCollection)
		// sort with score
		sort.SliceStable(tmpCollection, func(i, j int) bool {
			return vms[i] < vms[j]
		})

	} else {
		//don't need keys
		// have filter group
		if filter != nil {
			for i := 0; i < len(r.staticCollection); i++ {
				id := r.staticCollection[i]
				if !filter.Check(id) {
					tmpCollection = append(tmpCollection, id)
				}
			}
		} else {
			tmpCollection = r.staticCollection
		}
	}

	//TODO: Intersection with index accelerate

	// filter runtime condition
	if r.Condition != nil {
		tmpCollection = r.Condition.Check("user", userFeats, tmpCollection)
	}

	//order by
	if len(r.RecallEntityMeta.DSLMeta.OrderByMeta.Field) > 0 {
		poolsource.OrderBy(r.Pool, tmpCollection, r.RecallEntityMeta.DSLMeta.OrderByMeta.Field, r.RecallEntityMeta.DSLMeta.OrderByMeta.Desc)
	}

	//limit
	if r.RecallEntityMeta.Limit >= 0 {
		tmpCollection = tmpCollection[:r.RecallEntityMeta.Limit]
	}
	return tmpCollection
}

func (r *InvertInexRecall) Close() {
	if r.ref.CloseHandler != nil {
		r.ref.LazyFree(10)
	}
}
