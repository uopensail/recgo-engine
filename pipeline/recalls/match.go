package recalls

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
	"github.com/uopensail/ulib/sample"
)

type Matcher struct {
	conf    *model.MatchRecallConfigure
	index   *model.InvertedIndex
	program *minia.Minia
}

func NewMatcher(conf *model.MatchRecallConfigure) *Matcher {
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Expr)})
	return &Matcher{
		conf:    conf,
		index:   nil,
		program: program,
	}
}

func (m *Matcher) Do(uCtx *userctx.UserContext) model.Collection {
	var value sample.Features
	if uCtx.Related != nil {
		value = m.program.Eval(uCtx.Related, uCtx.Features)
	} else {
		value = m.program.Eval(uCtx.Features)
	}

	result := value.Get("result")

	if result == nil {
		return nil
	}

	keys, err := result.GetStrings()
	if err != nil {
		return nil
	}
	maxSize := 0
	for _, key := range keys {
		if entry, ok := m.index.Indexes[key]; ok {
			maxSize = max(maxSize, len(entry.Values))
		}
	}

	filter := make(map[int]struct{})
	ret := make([]*model.Entry, 0, m.conf.Count)
	for i := range maxSize {
		for _, key := range keys {
			if arr, ok := m.index.Indexes[key]; ok {
				if i < len(arr.Values) {
					k := arr.Values[i]
					if id, _ := uCtx.Items.GetByKey(k.Key); id >= 0 {
						if _, ok := filter[id]; ok {
							continue
						}
						filter[id] = struct{}{}
						entry, _ := model.NewEntry(k, uCtx.Items)
						entry.AddChan(m.conf.Name, fmt.Sprintf("recall by key: %s", key))
						ret = append(ret, entry)
					}
				}
			}
		}
	}

	count := min(m.conf.Count, len(ret))

	return ret[:count]
}
