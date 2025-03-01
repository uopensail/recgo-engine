package resources

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/finder"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/targz"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type Resource struct {
	table.PoolMeta
	*pool.Pool

	SubPoolCollectionRess SubPoolCollectionResource
	InvertIndexRess       map[string]InvertIndexFileResource
	LastUpdateTime        int64
}

func loadPoolResource(envCfg config.EnvConfig, targzFile string) (*Resource, error) {
	//清理旧的
	os.RemoveAll(filepath.Join(envCfg.WorkDir, "resources"))

	// unzip
	targz.Extract(targzFile, filepath.Join(envCfg.WorkDir, "resources"))

	ps := Resource{
		InvertIndexRess: make(map[string]InvertIndexFileResource),
	}
	// 解析meta
	err := table.LoadMeta(filepath.Join(envCfg.WorkDir, "resources", "pool.meta.json"), &ps.PoolMeta)
	if err != nil {
		return nil, err
	}

	//加载pool

	pl, err := pool.NewPool(filepath.Join(envCfg.WorkDir, "resources", "pool.json.txt"))
	if err != nil {
		return nil, err
	}
	ps.Pool = pl

	//加载物料子集合
	subPoolCollection, err := NewSubPoolCollectionResource(envCfg,
		filepath.Join(envCfg.WorkDir, "resources", "subpool.json.txt"))
	if err != nil {
		return nil, err
	}
	ps.SubPoolCollectionRess = *subPoolCollection

	//加载 invertIndex
	invertIndexRess := make(map[string]InvertIndexFileResource)
	invertIndexDir := path.Join(envCfg.WorkDir, "resource", "invert_index")

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
		invertIndex, err := NewInvertIndexFileResource(envCfg, fullPath)
		if err != nil {
			zlog.LOG.Warn("NewInvertIndexFileResource error", zap.Error(err))
			continue
		}
		invertIndexRess[fileName] = *invertIndex
	}

	ps.InvertIndexRess = invertIndexRess
	return &ps, nil

}
func NewResource(envCfg config.EnvConfig, remoteLocation string) (*Resource, error) {
	ps := Resource{}
	myFinder := finder.GetFinder(&envCfg.Finder)
	ps.LastUpdateTime = myFinder.GetUpdateTime(remoteLocation)
	localLoction := getFile(envCfg, remoteLocation)
	pl, err := loadPoolResource(envCfg, localLoction)
	if err != nil {
		zlog.LOG.Error("loadPoolResource", zap.Error(err))
		return nil, err
	}

	return pl, nil
}
func getFile(envCfg config.EnvConfig, location string) string {
	if strings.HasPrefix(location, "oss://") || strings.HasPrefix(location, "s3://") {
		baseName := filepath.Base(location)

		localPath := filepath.Join(envCfg.WorkDir, "tmp", baseName)
		myFinder := finder.GetFinder(&envCfg.Finder)
		myFinder.Download(location, localPath)
		return localPath
	} else {
		return location
	}

}

func (sm *Resource) CheckUpdateJob(newConf table.PoolMeta, envCfg config.EnvConfig) bool {
	oldConf := sm.PoolMeta
	//source meta 有更新
	needUpdate := false
	if oldConf.GetUpdateTime() != newConf.GetUpdateTime() {
		needUpdate = true
	} else {
		myFinder := finder.GetFinder(&envCfg.Finder)
		nUpdateTime := myFinder.GetUpdateTime(newConf.Location)
		if sm.LastUpdateTime < nUpdateTime {
			needUpdate = true
		}
	}
	return needUpdate

}
