package resources

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/uopensail/recgo-engine/config"
	xutils "github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"

	"go.uber.org/zap"
)

type ItemScore map[int]float32

func NewItemScore(location string, pl *pool.Pool) (ItemScore, error) {
	ret := make(map[int]float32)
	err := xutils.FileReadLine(location, func(line string) {
		vvs := strings.Split(line, "\t")
		if len(vvs) < 2 {
			return
		}

		key := vvs[0]
		v := utils.String2Float32(vvs[1])
		item := pl.GetByKey(key)
		if item == nil {
			return
		}
		ret[item.ID] = v
	})

	if err != nil {
		zlog.LOG.Error("failed to load file", zap.Error(err))
		return nil, err
	}
	return ret, nil
}

type ItemScoresResource struct {
	score map[string]ItemScore
}

func NewItemScoreResource(envCfg config.EnvConfig, location string, pl *pool.Pool) (ItemScoresResource, error) {
	resDir := path.Join(location, "scores")
	ress := ItemScoresResource{
		score: make(map[string]ItemScore),
	}
	entries, err := os.ReadDir(resDir)
	if err != nil {
		return ress, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(resDir, entry.Name())

		fileName := entry.Name()
		itemScore, err := NewItemScore(fullPath, pl)
		if err != nil {
			zlog.LOG.Warn("NewItemScore error", zap.Error(err))
			continue
		}
		ress.score[fileName] = itemScore
	}

	return ress, nil
}

type itemScoreHelper struct {
	Collection
	score []float32
}

func (x itemScoreHelper) Len() int { // 重写 Len() 方法
	return len(x.Collection)
}
func (x itemScoreHelper) Swap(i, j int) { // 重写 Swap() 方法
	x.Collection[i], x.Collection[j] = x.Collection[j], x.Collection[i]
	x.score[i], x.score[j] = x.score[j], x.score[i]
}
func (x itemScoreHelper) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return x.score[i] > x.score[j]
}

func (res *ItemScoresResource) SortByKey(key string, collection Collection) Collection {
	if itemScore, ok := res.score[key]; ok {
		scoreHelper := itemScoreHelper{
			Collection: make(Collection, len(collection)),
			score:      make([]float32, len(collection)),
		}
		copy(scoreHelper.Collection, collection)
		for i := 0; i < len(collection); i++ {
			score, has := itemScore[collection[i]]
			if has {
				scoreHelper.score[i] = score
			} else {
				scoreHelper.score[i] = -9999999
			}
		}
		sort.Sort(scoreHelper)
		return scoreHelper.Collection
	}
	return collection
}
