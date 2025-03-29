package strategy

import (
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/recgo-engine/strategy/freqfilter"
	"github.com/uopensail/recgo-engine/strategy/insert"
	"github.com/uopensail/recgo-engine/strategy/rank"
	"github.com/uopensail/recgo-engine/strategy/recalls/recall"
	"github.com/uopensail/recgo-engine/strategy/scatter"
	"github.com/uopensail/recgo-engine/strategy/weighted"
	"go.uber.org/zap"

	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/finder"
	xutils "github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
)

type Entities struct {
	ModelEntities
	Version int64
}

type EntitiesManager struct {
	entities *Entities
}

func (mgr *EntitiesManager) Init(envCfg config.EnvConfig, dataURL string, jobUtil *utils.MetuxJobUtil) {
	mgr.entities = &Entities{}

	mgr.cronJob(envCfg, dataURL, jobUtil)
}
func (mgr *EntitiesManager) GetEntities() *Entities {
	entities := (*Entities)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&mgr.entities))))
	return entities
}

func (mgr *EntitiesManager) cronJob(envCfg config.EnvConfig, dataURL string, jobUtil *utils.MetuxJobUtil) {
	job, err := mgr.loadAllJob(envCfg, dataURL)
	if err != nil {
		zlog.LOG.Error("loadAllJob", zap.Error(err))
		panic(err)

	}
	job()
	go func() {

		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()
		for {
			<-ticker.C
			job, err := mgr.loadAllJob(envCfg, dataURL)
			if err != nil {
				zlog.LOG.Error("loadAllJob", zap.Error(err))
				continue
			}
			jobUtil.TryRun(job)
		}
	}()

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

func (mgr *EntitiesManager) findLatestEngineDataDir(envCfg config.EnvConfig, locationPrefix string) (string, int64, error) {
	parentDir := filepath.Dir(locationPrefix)
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return "", -1, err
	}
	maxVersion := -1
	fileNamePrefix := filepath.Base(locationPrefix)
	maxDir := ""
	for _, entry := range entries {

		if entry.IsDir() {
			dirName := entry.Name()
			if (strings.HasPrefix(dirName, fileNamePrefix)) == false {
				continue
			}
			// 检查对应.success
			lockPath := filepath.Join(parentDir, dirName+".success")
			if _, err := os.Stat(lockPath); os.IsNotExist(err) {
				continue // 存在.success目录则跳过
			}
			vv := strings.Split(dirName, ".")
			if len(vv) != 2 {
				continue
			}
			version := xutils.String2Int(vv[1])

			// 比较版本号
			if version > maxVersion {
				maxVersion = version
				maxDir = filepath.Join(parentDir, entry.Name())
			}
		}
	}
	return maxDir, int64(maxVersion), nil

}

// Do not modify the execution order
func (mgr *EntitiesManager) loadAllJob(envCfg config.EnvConfig, url string) (func(), error) {

	engineDataDir, version, err := mgr.findLatestEngineDataDir(envCfg, url)
	if err != nil {
		zlog.LOG.Error("LoadDBTabelModel", zap.Error(err))
		return nil, err
	}

	oldEntities := mgr.entities

	if version <= oldEntities.Version {
		//不需要更新
		zlog.LOG.Info("Engine Data Don't need Update, use version:", zap.Int64("version", oldEntities.Version))
		return nil, nil
	}

	tableModel, err := dbmodel.LoadDBTabelModel(filepath.Join(engineDataDir, "dbmodel.toml"))

	if err != nil {
		zlog.LOG.Error("LoadDBTabelModel", zap.Error(err))
		return nil, err
	}

	sourceJobs := make([]func() bool, 0)
	entities := &Entities{
		ModelEntities: ModelEntities{
			Model: tableModel,
		},
	}
	if oldEntities != nil {
		entities.ModelEntities.Ress = oldEntities.Ress
	}
	resourceJob := func() bool {
		//TODO 优化Resource 变化了才更新
		ress, err := resources.NewResource(envCfg, filepath.Join(engineDataDir, "resources"))
		if err != nil {
			zlog.LOG.Error("NewResource", zap.Error(err))
			return false
		}
		entities.ModelEntities.Ress = *ress
		return true
	}
	sourceJobs = append(sourceJobs, resourceJob)
	entities.FilterResources.Clone(&oldEntities.FilterResources)
	job := entities.FilterResources.Reload(tableModel.FilterResourceTableModel.Rows, envCfg)
	if job != nil {
		sourceJobs = append(sourceJobs, job)
	}

	entityJob := func() bool {
		//Do not modify the execution order

		entities.FilterEntities = *freqfilter.NewFilterEntities(tableModel.FilterEntityTableModel.Rows, envCfg)

		entities.RecallEntities = *recall.NewRecallEntities(tableModel.RecallEntityTableModel.Rows, envCfg, &entities.Ress,
			&tableModel)

		entities.InsertEntities = *insert.NewInsertEntities(tableModel.InsertEntityTableModel.Rows, envCfg, &entities.Ress)

		entities.ScatterEntities = *scatter.NewScatterEntities(tableModel.ScatterEntityTableModel.Rows, envCfg)

		entities.RankEntities = *rank.NewRankEntities(tableModel.RankEntityTableModel.Rows, envCfg)
		entities.WeightedEntities = *weighted.NewWeightedEntities(tableModel.WeightedEntityTableModel.Rows, envCfg)
		entities.StrategyEntities = *NewStrategyEntities(tableModel.StrategyEntityTableModel.Rows, envCfg)
		entities.Model = tableModel
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&mgr.entities)), unsafe.Pointer(entities))
		return true
	}

	//如果source job 没有就立马更新
	if len(sourceJobs) == 0 {
		entityJob()
		return nil, nil
	} else {
		sourceJobs = append(sourceJobs, entityJob)
		return func() {
			for _, job := range sourceJobs {
				if job != nil {
					isContinue := job()
					if isContinue == false {
						zlog.LOG.Info("load job break", zap.Bool("isContinue", isContinue))
						return
					}
				}
			}
		}, nil
	}
}

var EntitiesMgr EntitiesManager
