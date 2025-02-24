package poolsource

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/finder"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type PoolSource struct {
	table.PoolMeta
	*pool.Pool
	*PoolIndeces
	LastUpdateTime int64
}

func NewPoolSource(newConf table.PoolMeta, envCfg config.EnvConfig) *PoolSource {
	ps := PoolSource{
		PoolMeta: newConf,
	}
	myFinder := finder.GetFinder(&envCfg.Finder)
	ps.LastUpdateTime = myFinder.GetUpdateTime(newConf.Location)
	pl, err := pool.NewPool(newConf.Location)
	if err != nil {
		zlog.LOG.Error("pool.NewPool", zap.Error(err))
		return nil
	}
	ps.Pool = pl
	ps.PoolIndeces = NewPoolIndeces(pl)
	return &ps
}

func (sm *PoolSource) CheckUpdateJob(newConf table.PoolMeta, envCfg config.EnvConfig) bool {
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
