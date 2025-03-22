package strategy

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/config"

	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type PluginCreateFunc func(cfg table.StrategyEntityMeta, env config.EnvConfig) IStrategyEntity

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

func PluginFactoryCreate(cfg table.StrategyEntityMeta, envCfg config.EnvConfig) IStrategyEntity {

	if createFunc, ok := pluginFactorys[cfg.PluginName]; ok {
		if createFunc != nil {
			return createFunc(cfg, envCfg)
		}
	}
	return nil
}

type StrategyEntities struct {
	Entities map[string]IStrategyEntity
}

func NewStrategyEntities(newConfs []table.StrategyEntityMeta, envCfg config.EnvConfig) *StrategyEntities {
	entities := &StrategyEntities{
		Entities: make(map[string]IStrategyEntity),
	}

	for _, v := range newConfs {
		s := PluginFactoryCreate(v, envCfg)
		if s != nil {
			entities.Entities[v.Name] = s
		}
	}
	return entities
}

func (entities *StrategyEntities) GetStrategy(name string) IStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("match.GetStrategy.%s", name))
	defer stat.End()

	if entity, ok := entities.Entities[name]; ok {
		return entity
	}
	stat.MarkErr()
	return nil

}
