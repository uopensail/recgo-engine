package resources

import (
	"sort"
	"strings"

	"github.com/uopensail/recgo-engine/config"
	xutils "github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/finder"
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
	copy(oc.Collection, oc.OrderCollection)
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

func NewInvertIndexFileResource(envCfg config.EnvConfig, location string) (*InvertIndexFileResource, error) {

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
		vv := utils.String2IntList(vvs[1], ",")
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

func (res *InvertIndexFileResource) CheckResourceUpdate(envCfg config.EnvConfig, poolUpdate bool) bool {
	if poolUpdate {
		return true
	}
	myFinder := finder.GetFinder(&envCfg.Finder)
	nUpdateTime := myFinder.GetUpdateTime(res.location)
	if res.updateTime < nUpdateTime {
		return true
	}
	return false
}

func (res *InvertIndexFileResource) Close() {

}
