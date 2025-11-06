package model

import (
	"fmt"

	"github.com/uopensail/ulib/sample"
)

type Runtime struct {
	Basic   *sample.ImmutableFeatures
	RunTime *sample.MutableFeatures
}

func NewRuntime(basic *sample.ImmutableFeatures) *Runtime {
	return &Runtime{
		Basic:   basic,
		RunTime: sample.NewMutableFeatures(),
	}
}

func (r *Runtime) Get(key string) (sample.Feature, error) {
	if fea := r.Basic.Get(key); fea != nil {
		return fea, nil
	}

	fea := r.RunTime.Get(key)
	if fea != nil {
		return fea, nil
	}
	return nil, fmt.Errorf("key: %s miss", key)
}

func (r *Runtime) Set(key string, value sample.Feature) {
	r.RunTime.Set(key, value)
}
