package resource

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
)

func Create(cfg table.FilterResourceMeta, env config.EnvConfig) Resource {
	switch cfg.PluginName {
	case "redis":
		return NewRedisResource(cfg, env)
	}
	return NewRedisResource(cfg, env)
}
