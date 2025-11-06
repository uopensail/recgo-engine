package constrains

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
)

type FixedPositionInsert struct {
	conf    *model.FixedPositionInsertedConstrainConfigure
	program *minia.Minia
}

func NewFixedPositionInsert(conf *model.FixedPositionInsertedConstrainConfigure) *FixedPositionInsert {
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Condition)})
	return &FixedPositionInsert{
		conf:    conf,
		program: program,
	}
}

func (f *FixedPositionInsert) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	if f.conf.Position >= len(collection) {
		return collection
	}

	ret := make([]*model.Entry, 0, len(collection))
	for i, entry := range collection {

		value := f.program.Eval(entry.Runtime.Basic, entry.Runtime.RunTime, uCtx.Features)
		result := value.Get("result")
		if result == nil {
			continue
		}

		hit, err := result.GetInt64()
		if err == nil && hit == 1 {
			if i == f.conf.Position {
				return collection
			} else if i < f.conf.Position {
				slice1 := collection[:i]
				slice2 := collection[i+1 : f.conf.Position]
				slice3 := collection[f.conf.Position:]
				ret = append(ret, slice1...)
				ret = append(ret, slice2...)
				ret = append(ret, entry)
				ret = append(ret, slice3...)
				return ret
			} else {
				slice1 := collection[:f.conf.Position]
				slice2 := collection[f.conf.Position:i]
				slice3 := collection[i+1:]
				ret = append(ret, slice1...)
				ret = append(ret, entry)
				ret = append(ret, slice2...)
				ret = append(ret, slice3...)
				return ret
			}
		}
	}
	return collection
}
