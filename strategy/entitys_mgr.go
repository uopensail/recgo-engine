package strategy

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"go.uber.org/zap"

	"github.com/uopensail/recgo-engine/poolsource"

	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/zlog"
)

type EntitiesManager struct {
	entities *ModelEntities
}

func (mgr *EntitiesManager) Init(envCfg config.EnvConfig, jobUtil *utils.MetuxJobUtil) {
	mgr.entities = &ModelEntities{}

	mgr.cronJob(envCfg, jobUtil)
}
func (mgr *EntitiesManager) GetModelEntities() *ModelEntities {
	entities := (*ModelEntities)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&mgr.entities))))
	return entities
}

func (mgr *EntitiesManager) cronJob(envCfg config.EnvConfig, jobUtil *utils.MetuxJobUtil) {
	job := mgr.loadAllJob(envCfg)
	job()
	go func() {

		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()
		for {
			<-ticker.C
			job := mgr.loadAllJob(envCfg)
			jobUtil.TryRun(job)
		}
	}()

}

// Do not modify the execution order
func (mgr *EntitiesManager) loadAllJob(envCfg config.EnvConfig) func() {
	tableModel, err := dbmodel.LoadDBTabelModel(config.AppConfigInstance.URL)
	if err != nil {
		zlog.LOG.Error("LoadDBTabelModel", zap.Error(err))
		return nil
	}
	oldEntities := mgr.entities

	entities := &ModelEntities{
		PoolSource: oldEntities.PoolSource,
	}
	sourceJobs := make([]func(), 0)
	var poolUpdate bool
	if len(tableModel.PoolSourceTableModel.Rows) > 0 {

		if entities.PoolSource.CheckUpdateJob(tableModel.PoolSourceTableModel.Rows[0], envCfg) {
			sourceJobs = append(sourceJobs, func() {
				ps := poolsource.NewPoolSource(tableModel.PoolSourceTableModel.Rows[0], envCfg)
				if ps != nil {
					entities.PoolSource = *ps

					*(&poolUpdate) = true
				}
			})

		}

	}
	entities.RecallResources.Clone(&oldEntities.RecallResources)
	job := entities.RecallResources.Reload(envCfg, tableModel.RecallSourceTableModel.Rows, &entities.PoolSource, poolUpdate)
	if job != nil {
		sourceJobs = append(sourceJobs, job)
	}

	entities.FilterResources.Clone(&oldEntities.FilterResources)
	job = entities.FilterResources.Reload(tableModel.FilterResourceTableModel.Rows, envCfg)
	if job != nil {
		sourceJobs = append(sourceJobs, job)
	}

	entityJob := func() {
		//Do not modify the execution order
		entities.FilterEntities.Clone(&oldEntities.FilterEntities)
		entities.FilterEntities.Reload(tableModel.FilterEntityTableModel.Rows, envCfg)

		entities.RecallEntities.Clone(&oldEntities.RecallEntities)
		entities.RecallEntities.Reload(tableModel.RecallEntityTableModel.Rows, envCfg,
			entities.PoolSource.Pool, poolUpdate, &tableModel)
		entities.InsertEntities.Clone(&oldEntities.InsertEntities)
		entities.InsertEntities.Reload(tableModel.InsertEntityTableModel.Rows, envCfg, entities.PoolSource.Pool)

		entities.ScatterEntities.Clone(&oldEntities.ScatterEntities)
		entities.ScatterEntities.Reload(tableModel.ScatterEntityTableModel.Rows, envCfg)
		entities.ScatterEntities.Clone(&oldEntities.ScatterEntities)
		entities.ScatterEntities.Reload(tableModel.ScatterEntityTableModel.Rows, envCfg)

		entities.RankEntities.Clone(&oldEntities.RankEntities)
		entities.RankEntities.Reload(tableModel.RankEntityTableModel.Rows, envCfg)
		entities.WeightedEntities.Clone(&oldEntities.WeightedEntities)
		entities.WeightedEntities.Reload(tableModel.WeightedEntityTableModel.Rows, envCfg)
		entities.StrategyEntities.Clone(&oldEntities.StrategyEntities)
		entities.StrategyEntities.Reload(tableModel.StrategyEntityTableModel.Rows, envCfg)
		entities.Model = tableModel
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&mgr.entities)), unsafe.Pointer(entities))
	}

	//如果source job 没有就立马更新
	if len(sourceJobs) == 0 {
		entityJob()
		return nil
	} else {
		sourceJobs = append(sourceJobs, entityJob)
		return func() {
			for _, job := range sourceJobs {
				if job != nil {
					job()
				}
			}
		}
	}
}

var EntitiesMgr EntitiesManager
