package freqfilter

import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/strategy/freqfilter/resource"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/uno"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type FilterEntity struct {
	table.FilterEntityMeta
	unoRef    utils.Reference
	evaluator *uno.Evaluator
}

func NewFilterEntity(cfg table.FilterEntityMeta) *FilterEntity {
	stat := prome.NewStat("FilterResource.NewFilterResource")
	defer stat.End()
	var evaluator *uno.Evaluator
	if len(cfg.Condition) > 0 {
		var err error
		evaluator, err = uno.NewEvaluator(cfg.Condition)
		if err != nil {
			stat.MarkErr()
			zlog.LOG.Error("condition parse error", zap.Error(err))
			return nil
		}
	}

	entity := &FilterEntity{
		FilterEntityMeta: cfg,
		evaluator:        evaluator,
	}
	entity.unoRef.CloseHandler = func() {
		if entity.evaluator != nil {
			entity.evaluator.Release()
		}
	}
	return entity
}
func (r *FilterEntity) Meta() *table.FilterEntityMeta {
	return &r.FilterEntityMeta
}

func (r *FilterEntity) Close() {
	if r.unoRef.CloseHandler != nil {
		r.unoRef.LazyFree(10)
	}
}

func (r *FilterEntity) Do(useId string, ress *resource.Resources) []string {
	stat := prome.NewStat("FilterResource.Get")
	defer stat.End()
	iResource := ress.GetResource(r.FilterEntityMeta.SourceID)

	keyFormat := useId
	if len(r.FilterEntityMeta.Format) > 0 {
		keyFormat = fmt.Sprintf(r.FilterEntityMeta.Format, useId)
	}
	list, err := iResource.Do(context.TODO(), keyFormat, r.FilterEntityMeta.Params)

	if err != nil {
		stat.MarkErr()
		zlog.LOG.Error("redis lrange error", zap.Error(err))
		return nil
	}

	// data format: id|ts|scane|action,...
	feas := make([]sample.Features, 0, 1024)
	ids := make([]string, 0, 1024)
	for i := 0; i < len(list); i++ {
		items := strings.Split(list[i], ",")
		for j := 0; j < len(items); j++ {
			values := strings.Split(items[j], "|")
			features := sample.NewMutableFeatures()
			features.Set("id", &sample.String{Value: values[0]})
			features.Set("ts", &sample.Int64{Value: utils.String2Int64(values[1])})
			features.Set("scene", &sample.String{Value: values[2]})
			features.Set("action", &sample.String{Value: values[3]})
			feas = append(feas, features)
			ids = append(ids, values[0])
		}
	}

	// do condition filter
	var evalStatus []int32
	filtered := make(map[string]int, len(ids))
	if r.evaluator != nil {
		slice := r.evaluator.Allocate()
		slices := make([][]unsafe.Pointer, len(feas))
		for i := 0; i < len(feas); i++ {
			tmp := make([]unsafe.Pointer, len(slice))
			copy(tmp, slice)
			r.evaluator.Fill("", feas[i], tmp)
			slices[i] = tmp
		}
		evalStatus = r.evaluator.BatchEval(slices)
		for i := 0; i < len(ids); i++ {
			if len(evalStatus) > i && evalStatus[i] == 1 {
				if count, ok := filtered[ids[i]]; ok {
					filtered[ids[i]] = count + 1
				} else {
					filtered[ids[i]] = 1
				}
			}
		}
	} else {
		for i := 0; i < len(ids); i++ {
			if count, ok := filtered[ids[i]]; ok {
				filtered[ids[i]] = count + 1
			} else {
				filtered[ids[i]] = 1
			}
		}
	}

	// do aggregation filter
	ret := make([]string, 0, len(filtered))
	for id, count := range filtered {
		if count >= r.FilterEntityMeta.MaxCount {
			ret = append(ret, id)
		}
	}
	return ret
}
