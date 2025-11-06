package constrains

import (
	"strconv"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/sample"
)

type Scatter struct {
	confs []*model.ScatterBasedConstrainConfigure
}

func NewScatter(confs []*model.ScatterBasedConstrainConfigure) *Scatter {
	return &Scatter{
		confs: confs,
	}
}

func (s *Scatter) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	filter := make([]map[string]int, len(s.confs))
	for i := range s.confs {
		filter[i] = make(map[string]int)
	}

	ret := make([]*model.Entry, 0, len(collection))
	remain := make([]*model.Entry, 0, len(collection))

	keys := make([][]string, len(s.confs))
	for _, entry := range collection {
		status := true

		for i := 0; i < len(s.confs) && status; i++ {
			conf := s.confs[i]
			fea, err := entry.Get(conf.Field)
			if err == nil {
				keys[i] = nil
				continue
			}
			keys[i] = Feature2StringSlice(fea)

			for _, key := range keys[i] {
				if val, ok := filter[i][key]; ok {
					if val >= conf.Count {
						status = false
						break
					}
				}
			}
		}

		if status {
			ret = append(ret, entry)
			for i := 0; i < len(s.confs); i++ {
				for _, key := range keys[i] {
					if val, ok := filter[i][key]; ok {
						filter[i][key] = val + 1
					} else {
						filter[i][key] = 1
					}
				}
			}
		} else {
			remain = append(remain, entry)
		}
	}
	ret = append(ret, remain...)
	return ret
}

func Feature2StringSlice(feature sample.Feature) []string {
	switch feature.Type() {
	case sample.Int64Type:
		val, _ := feature.GetInt64()
		return []string{strconv.FormatInt(val, 10)}
	case sample.Int64sType:
		val, _ := feature.GetInt64s()
		ret := make([]string, 0, len(val))
		for _, v := range val {
			ret = append(ret, strconv.FormatInt(v, 10))
		}
		return ret
	case sample.StringType:
		val, _ := feature.GetString()
		return []string{val}
	case sample.StringsType:
		val, _ := feature.GetStrings()
		return val
	default:
		return nil
	}
}
