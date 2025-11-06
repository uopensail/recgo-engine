package constrains

import (
	"fmt"
	"sort"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
)

type WeightAdjust struct {
	conf    *model.WeightAdjustedConstrainConfigure
	program *minia.Minia
}

func NewWeightAdjust(conf *model.WeightAdjustedConstrainConfigure) *WeightAdjust {
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Condition)})
	return &WeightAdjust{
		conf:    conf,
		program: program,
	}
}

func (w *WeightAdjust) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	for _, entry := range collection {

		value := w.program.Eval(entry.Runtime.Basic, entry.Runtime.RunTime, uCtx.Features)
		if value != nil {
			result := value.Get("result")
			if result == nil {
				continue
			}
			hit, err := result.GetInt64()
			if err == nil && hit == 1 {
				entry.KeyScore.Score *= w.conf.Ratio
			}
		}
	}

	sort.Stable(collection)
	return collection
}
