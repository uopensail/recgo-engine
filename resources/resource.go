package resources

import (
	"os"
	"path"
	"path/filepath"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type ResourceMeta struct {
	FieldDataType map[string]sample.DataType `json:"field_data_type" toml:"field_data_type"`
}

type Resource struct {
	UserFieldDataType map[string]sample.DataType
	ResourceMeta
	*pool.Pool

	SubPoolCollectionRess SubPoolCollectionResource
	InvertIndexRess       map[string]InvertIndexFileResource
	ItemScoresResource
}

func loadPoolResource(envCfg config.EnvConfig, resourcesDir string) (*Resource, error) {

	ps := Resource{
		InvertIndexRess: make(map[string]InvertIndexFileResource),
	}
	// 解析meta
	err := table.LoadMeta(filepath.Join(resourcesDir, "resource.meta.json"), &ps.ResourceMeta)
	if err != nil {
		return nil, err
	}

	//加载pool

	pl, err := pool.NewPool(filepath.Join(resourcesDir, "pool.txt"))
	if err != nil {
		return nil, err
	}
	ps.Pool = pl

	//加载物料子集合
	subPoolCollection, err := NewSubPoolCollectionResource(envCfg,
		filepath.Join(resourcesDir, "subpool.txt"), pl)
	if err != nil {
		return nil, err
	}
	ps.SubPoolCollectionRess = *subPoolCollection

	//加载 invertIndex
	invertIndexRess := make(map[string]InvertIndexFileResource)
	invertIndexDir := path.Join(resourcesDir, "invert_index")

	entries, err := os.ReadDir(invertIndexDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(invertIndexDir, entry.Name())

		fileName := entry.Name()
		invertIndex, err := NewInvertIndexFileResource(envCfg, fullPath, pl)
		if err != nil {
			zlog.LOG.Warn("NewInvertIndexFileResource error", zap.Error(err))
			continue
		}
		invertIndexRess[fileName] = *invertIndex
	}

	ps.InvertIndexRess = invertIndexRess

	itemScoresResource, err := NewItemScoreResource(envCfg, resourcesDir, pl)
	if err != nil {
		zlog.LOG.Warn("NewItemScoreResource error", zap.Error(err))
	}
	ps.ItemScoresResource = itemScoresResource
	return &ps, nil

}
func NewResource(envCfg config.EnvConfig, localDir string) (*Resource, error) {

	ps, err := loadPoolResource(envCfg, localDir)
	if err != nil {
		zlog.LOG.Error("loadPoolResource", zap.Error(err))
		return nil, err
	}

	return ps, nil
}
