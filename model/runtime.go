package model

import (
	"fmt"

	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Runtime holds both immutable and mutable feature sets for an item.
// Immutable features represent static attributes loaded from storage.
// Mutable features are dynamically added/updated during runtime.
type Runtime struct {
	Basic   *sample.ImmutableFeatures // Immutable/static features
	RunTime *sample.MutableFeatures   // Mutable/dynamic features
}

// NewRuntime creates a new Runtime instance given immutable basic features.
// It initializes an empty set of mutable features.
func NewRuntime(basic *sample.ImmutableFeatures) *Runtime {
	return &Runtime{
		Basic:   basic,
		RunTime: sample.NewMutableFeatures(),
	}
}

// Get retrieves a feature by key.
// It first looks into immutable (Basic) features, then mutable (RunTime) features.
// Returns an error if the key is not found in both.
// Logs will record missing keys for easier debugging.
func (r *Runtime) Get(key string) (sample.Feature, error) {
	if fea := r.Basic.Get(key); fea != nil {
		return fea, nil
	}
	if fea := r.RunTime.Get(key); fea != nil {
		return fea, nil
	}

	zlog.LOG.Warn("Runtime.Get.KeyMissing", zap.String("key", key))
	return nil, fmt.Errorf("key: %s miss", key)
}

// Set adds or updates a mutable (runtime) feature value for the given key.
// Logs every set operation for tracing runtime modifications.
func (r *Runtime) Set(key string, value sample.Feature) {
	r.RunTime.Set(key, value)
	zlog.LOG.Info("Runtime.Set",
		zap.String("key", key),
		zap.Any("value", value),
	)
}
