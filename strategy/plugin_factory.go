package strategy

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/utils"

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
	Entities map[int]IStrategyEntity
}

func (Entities *StrategyEntities) Clone(a *StrategyEntities) {
	Entities.Entities = make(map[int]IStrategyEntity)
	if a != nil {
		for k, v := range a.Entities {
			Entities.Entities[k] = v
		}
	}
}

func (Entities *StrategyEntities) GetStrategy(id int) IStrategyEntity {
	stat := prome.NewStat(fmt.Sprintf("match.GetStrategy.%d", id))
	defer stat.End()

	if entity, ok := Entities.Entities[id]; ok {
		return entity
	}
	stat.MarkErr()
	return nil

}

func (Entities *StrategyEntities) Reload(newConfs []table.StrategyEntityMeta, envCfg config.EnvConfig) {
	oldConfs := make([]table.StrategyEntityMeta, 0, len(Entities.Entities))
	for _, v := range Entities.Entities {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)

	if len(invalidM)+len(upsertM) <= 0 {
		return
	}

	for k, v := range upsertM {
		s := PluginFactoryCreate(v, envCfg)
		if s != nil {
			Entities.Entities[k] = s
		}
	}

	//删除
	for k := range invalidM {
		delete(Entities.Entities, k)
	}

}
