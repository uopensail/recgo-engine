package model

import (
	"fmt"

	"github.com/uopensail/ulib/sample"
)

type IFliter interface {
	Check(id int) bool //True: pass, False: filted
	Exclude() []string
}

type Resource interface {
	GetUpdateTime() (int64, error)
	GetURL() (string, error)
}

const ChannelsKey = "i_ctx_chans"
const ReasonsKey = "i_ctx_reasons"

type Entry struct {
	ID int
	KeyScore
	Runtime
}

func NewEntry(k KeyScore, items *Items) (*Entry, error) {
	id, feas := items.GetByKey(k.Key)
	if feas == nil {
		return nil, fmt.Errorf("key miss: %s", k.Key)
	}

	r := NewRuntime(feas)
	r.Set(ChannelsKey, &sample.Strings{Value: make([]string, 0, 8)})
	r.Set(ReasonsKey, &sample.Strings{Value: make([]string, 0, 8)})
	return &Entry{id, k, *r}, nil
}

func (entry *Entry) AddChan(ch string, reason string) {
	var fea sample.Feature
	fea, _ = entry.Get(ChannelsKey)

	tmp, _ := fea.GetStrings()
	tmp = append(tmp, ch)

	fea, _ = entry.Get(ReasonsKey)
	tmp, _ = fea.GetStrings()
	tmp = append(tmp, reason)
}

func (entry *Entry) MergeChans(e *Entry) {
	f0, _ := entry.Get(ChannelsKey)
	f1, _ := e.Get(ChannelsKey)

	tmp0, _ := f0.GetStrings()
	tmp1, _ := f1.GetStrings()

	tmp0 = append(tmp0, tmp1...)

	f0, _ = entry.Get(ReasonsKey)
	f1, _ = e.Get(ReasonsKey)
	tmp0, _ = f0.GetStrings()
	tmp1, _ = f1.GetStrings()

	tmp0 = append(tmp0, tmp1...)
}

type Collection []*Entry

func (c Collection) Less(i, j int) bool {
	return c[i].Score > c[j].Score
}

func (c Collection) Len() int {
	return len(c)
}

func (list Collection) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
