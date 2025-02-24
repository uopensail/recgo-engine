package resource

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/pool"
)

func Create(envCfg config.EnvConfig, meta table.RecallResourceMeta, pl *pool.Pool) IResource {
	switch meta.PluginName {
	case "file":
		return NewFileResource(envCfg, meta, pl)
	case "redis":
		return NewRedisResource(envCfg, meta, pl)
	default:
		return nil
	}

}
