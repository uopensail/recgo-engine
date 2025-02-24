package insert

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type PluginCreateFunc func(cfg table.InsertEntityMeta, env config.EnvConfig, pl *pool.Pool) IStrategyEntity

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

func PluginFactoryCreate(cfg table.InsertEntityMeta, envCfg config.EnvConfig, pl *pool.Pool) IStrategyEntity {

	if createFunc, ok := pluginFactorys[cfg.PluginName]; ok {
		if createFunc != nil {
			return createFunc(cfg, envCfg, pl)
		}
	}
	return nil
}
