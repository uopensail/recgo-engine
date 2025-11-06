package resources

import (
	"github.com/uopensail/recgo-engine/model"

	"github.com/uopensail/recgo-engine/config"
)

// Default check interval in seconds
const interval int = 300

type ResourceManager struct {
	indexes map[string]*Finder
	items   *Finder
}

func NewResourceManager(conf *config.AppConfig) *ResourceManager {
	indexes := make(map[string]*Finder, len(conf.Indexes))
	for _, res := range conf.Indexes {
		index, err := NewFinder(res.Dir, model.NewInvertedIndex)
		if err != nil {
			panic(err)
		}
		indexes[res.Name] = index
	}

	items, err := NewFinder(conf.Items.Dir, model.NewItems)
	if err != nil {
		panic(err)
	}
	return &ResourceManager{indexes: indexes, items: items}
}

func (m *ResourceManager) GetItems() *model.Items {
	res := m.items.Get()
	return res.(*model.Items)
}

func (m *ResourceManager) GetIndex(name string) *model.InvertedIndex {
	if index, ok := m.indexes[name]; ok {
		res := index.Get()
		return res.(*model.InvertedIndex)
	}
	return nil
}

var ResourceManagerInstance *ResourceManager
