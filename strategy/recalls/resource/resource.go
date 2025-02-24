package resource

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/poolsource"
	"github.com/uopensail/recgo-engine/utils"
	"github.com/uopensail/ulib/datastruct"
	"github.com/uopensail/ulib/pool"
)

type IResource interface {
	Get(keys []string, pl *pool.Pool) [][]datastruct.Tuple[int, float32]
	CheckResourceUpdate(envCfg config.EnvConfig, poolUpdate bool) bool
	Meta() *table.RecallResourceMeta
	Close()
}

type Resources struct {
	resources map[string]IResource
}

func NewResources(envCfg config.EnvConfig, metas []table.RecallResourceMeta, pl *pool.Pool) *Resources {
	ress := Resources{
		resources: make(map[string]IResource, len(metas)),
	}
	for _, meta := range metas {
		res := Create(envCfg, meta, pl)
		if res != nil {
			ress.resources[meta.Name] = res
		}
	}
	return &ress
}

func (ress *Resources) Get(name string) IResource {
	if v, ok := ress.resources[name]; ok {
		return v
	}
	return nil
}

func (ress *Resources) Close() {
	for _, v := range ress.resources {
		v.Close()
	}
}

func (sm *Resources) Clone(a *Resources) {
	sm.resources = make(map[string]IResource)
	if a != nil {
		for k, v := range a.resources {
			sm.resources[k] = v
		}
	}
}

func (sm *Resources) Reload(envCfg config.EnvConfig, newConfs []table.RecallResourceMeta, plSource *poolsource.PoolSource, poolUpdate bool) func() {
	oldConfs := make([]table.RecallResourceMeta, 0, len(sm.resources))
	for _, v := range sm.resources {
		cfg := v.Meta()
		oldConfs = append(oldConfs, *cfg)
	}

	//source meta 有更新
	invalidM, upsertM := utils.CheckUpsert(oldConfs, newConfs)
	// 检查soource 的实际内容是否有更新

	for _, v := range sm.resources {
		vMeta := v.Meta()
		if _, ok := invalidM[vMeta.ID]; ok {
			//即将要删除的
			continue
		}
		if _, ok := upsertM[vMeta.ID]; ok {
			//即将要重新new的
			continue
		}
		needUpdate := v.CheckResourceUpdate(envCfg, poolUpdate)
		if needUpdate {
			upsertM[vMeta.ID] = *vMeta
		}
	}

	if len(invalidM)+len(upsertM) <= 0 {
		return nil
	}

	return func() {

		for _, v := range upsertM {
			id := v.Name
			s := Create(envCfg, v, plSource.Pool)
			if s != nil {
				if old, ok := sm.resources[id]; ok {
					old.Close()
				}
				sm.resources[id] = s
			}
		}

		//删除
		for k := range invalidM {
			for name, v := range sm.resources {
				if v.Meta().ID == k {
					v.Close()
					delete(sm.resources, name)
					break
				}
			}
		}

	}

}
