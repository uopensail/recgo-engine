package resources

import (
	"sort"
	"strings"

	"github.com/uopensail/recgo-engine/config"
	xutils "github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"

	"go.uber.org/zap"
)

type OrderCollection struct {
	Collection                 //无序原始集合
	OrderCollection Collection // 有序集合，用于判断是否存在，求交集等
}

func NewOrderCollection(vv []int) OrderCollection {
	oc := OrderCollection{
		Collection: vv,
	}
	oc.OrderCollection = make(Collection, len(vv))
	copy(oc.OrderCollection, oc.Collection)
	sort.Ints(oc.OrderCollection)
	return oc

}

type fileIndeces struct {
	updateTime int64
	indexes    map[string]OrderCollection //索引无序集合cache
}

type InvertIndexFileResource struct {
	location string
	*fileIndeces
}

func NewInvertIndexFileResource(envCfg config.EnvConfig, location string, pl *pool.Pool) (*InvertIndexFileResource, error) {

	fs := &InvertIndexFileResource{
		location: location,
		fileIndeces: &fileIndeces{
			indexes: make(map[string]OrderCollection),
		},
	}
	indeces := fileIndeces{
		indexes: make(map[string]OrderCollection),
	}
	err := xutils.FileReadLine(location, func(line string) {
		vvs := strings.Split(line, "\t")
		if len(vvs) < 2 {
			return
		}

		key := vvs[0]
		vvStr := utils.StringSplit(vvs[1], ",")
		vv := make([]int, 0, len(vvStr))
		for _, v := range vvStr {
			item := pl.GetByKey(v)
			if item == nil {
				continue
			}
			vv = append(vv, item.ID)
		}

		indeces.indexes[key] = NewOrderCollection(vv)
	})
	fs.fileIndeces = &indeces
	if err != nil {
		zlog.LOG.Error("failed to load file", zap.Error(err))
		return nil, err
	}
	return fs, nil
}

func (res *InvertIndexFileResource) Get(keys []string) []Collection {
	ret := make([]Collection, len(keys))
	for i := 0; i < len(keys); i++ {
		v := res.indexes[keys[i]]

		ret[i] = make([]int, len(v.Collection))
		copy(ret[i], v.Collection)
	}
	return ret
}

func (res *InvertIndexFileResource) Close() {

}
