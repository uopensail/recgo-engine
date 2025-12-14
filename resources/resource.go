package resources

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// ResourceManager manages Items and Index resources, periodically reloading them.
type ResourceManager struct {
	indexes map[string]*Finder
	items   *Finder
}

// NewResourceManager creates and initializes all resources defined in the AppConfig.
//
// @param conf Application configuration.
// @return ResourceManager pointer.
func NewResourceManager(conf *config.AppConfig) *ResourceManager {
	pStat := prome.NewStat("NewResourceManager")
	defer pStat.End()

	indexes := make(map[string]*Finder, len(conf.Indexes))

	// Initialize each index finder
	for _, res := range conf.Indexes {
		index, err := NewFinder(res.Dir, model.NewInvertedIndex)
		if err != nil {
			zlog.LOG.Fatal("ResourceManager: failed to initialize index",
				zap.String("name", res.Name),
				zap.String("dir", res.Dir),
				zap.Error(err))
		}
		indexes[res.Name] = index
	}

	// Initialize items finder
	items, err := NewFinder(conf.Items.Dir, model.NewItems)
	if err != nil {
		zlog.LOG.Fatal("ResourceManager: failed to initialize items",
			zap.String("dir", conf.Items.Dir),
			zap.Error(err))
	}

	rm := &ResourceManager{
		indexes: indexes,
		items:   items,
	}

	zlog.LOG.Info("ResourceManager: initialized successfully",
		zap.Int("indexes_count", len(indexes)),
		zap.String("items_dir", conf.Items.Dir))
	return rm
}

// GetItems returns the current Items resource.
func (m *ResourceManager) GetItems() *model.Items {
	res := m.items.Get()
	return res.(*model.Items)
}

// GetIndex returns the current InvertedIndex resource by name.
// Returns nil if index is not found.
func (m *ResourceManager) GetIndex(name string) *model.InvertedIndex {
	if index, ok := m.indexes[name]; ok {
		res := index.Get()
		return res.(*model.InvertedIndex)
	}
	return nil
}

// ResourceManagerInstance is the global singleton instance.
var ResourceManagerInstance *ResourceManager
