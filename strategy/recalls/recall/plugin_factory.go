package recall

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

func init() {
	RegisterPlugin("invert_index", NewInvertInexRecall)
}

type IRecallStrategyEntity interface {
	Do(uCtx *userctx.UserContext, ifilter model.IFliter) ([]int, error)
	Meta() *table.RecallEntityMeta
	Close()
}
type PluginCreateFunc func(cfg table.RecallEntityMeta, pl *pool.Pool, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity

var pluginFactorys map[string]PluginCreateFunc

func RegisterPlugin(name string, createFunc PluginCreateFunc) {
	if pluginFactorys == nil {
		pluginFactorys = make(map[string]PluginCreateFunc)
	}
	if _, ok := pluginFactorys[name]; ok == false {
		pluginFactorys[name] = createFunc
	} else {
		zlog.LOG.Error("source plugin already exists", zap.String("name", name))
	}
}

func PluginFactoryCreate(cfg table.RecallEntityMeta, pl *pool.Pool, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity {

	if createFunc, ok := pluginFactorys[cfg.PluginName]; ok {
		if createFunc != nil {
			return createFunc(cfg, pl, dbModel)
		}
	}
	return nil
}
