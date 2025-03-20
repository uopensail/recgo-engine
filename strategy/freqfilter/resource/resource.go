package resource

import (
	"context"
	"fmt"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/prome"
)

// 定义数据源
type Resource interface {
	Meta() *table.FilterResourceMeta
	Do(ctx context.Context, key string, param map[string]string) ([]string, error)
	Close()
}

type Resources struct {
	sources map[int]Resource
}

func (sm *Resources) Clone(a *Resources) {
	sm.sources = make(map[int]Resource)
	if a != nil {
		for k, v := range a.sources {
			sm.sources[k] = v
		}
	}
}

func (sm *Resources) GetResource(id int) Resource {
	stat := prome.NewStat(fmt.Sprintf("recall.Source.%d", id))
	defer stat.End()

	if s, ok := sm.sources[id]; ok {
		return s
	}
	stat.MarkErr()
	return nil
}

func (sm *Resources) Reload(newConfs []table.FilterResourceMeta, envCfg config.EnvConfig) func() bool {
	oldConfs := make([]table.FilterResourceMeta, 0, len(sm.sources))
	for _, v := range sm.sources {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	//source meta 有更新
	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)

	if len(invalidM)+len(upsertM) <= 0 {
		return nil
	}

	return func() bool {

		for k, v := range upsertM {
			s := Create(v, envCfg)
			if s != nil {
				if old, ok := sm.sources[k]; ok {
					old.Close()
				}
				sm.sources[k] = s
			}
		}

		//删除
		for k := range invalidM {
			if old, ok := sm.sources[k]; ok {
				old.Close()
			}
			delete(sm.sources, k)
		}
		return true
	}

}
