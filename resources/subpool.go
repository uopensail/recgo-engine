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

type SubPoolCollectionResource struct {
	SubPool map[int]Collection // 有序集合
}

func NewSubPoolCollectionResource(envCfg config.EnvConfig, location string, pl *pool.Pool) (*SubPoolCollectionResource, error) {

	fs := &SubPoolCollectionResource{
		SubPool: make(map[int]Collection),
	}
	// 加载物料子集合
	subPoolCollection := make(map[int]Collection)
	err := xutils.FileReadLine(location, func(line string) {
		vvs := strings.Split(line, "\t")
		if len(vvs) < 2 {
			return
		}

		key := utils.String2Int(vvs[0])
		vvStr := utils.StringSplit(vvs[1], ",")
		vv := make([]int, 0, len(vvStr))
		for _, v := range vvStr {
			item := pl.GetByKey(v)
			vv = append(vv, item.ID)
		}
		sort.Ints(vv)
		subPoolCollection[key] = vv
	})
	fs.SubPool = subPoolCollection
	if err != nil {
		zlog.LOG.Error("failed to load file", zap.Error(err))
		return nil, err
	}
	return fs, nil
}

func (res *SubPoolCollectionResource) Get(key int) Collection {
	if v, ok := res.SubPool[key]; ok {
		return v
	}
	return nil
}

func (res *SubPoolCollectionResource) CheckResourceUpdate(envCfg config.EnvConfig, poolUpdate bool) bool {

	return false
}

func (res *SubPoolCollectionResource) Close() {

}
