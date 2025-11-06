package pipeline

import (
	"fmt"
	"sync"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/pipeline/constrains"
	"github.com/uopensail/recgo-engine/pipeline/freqs"
	"github.com/uopensail/recgo-engine/pipeline/rank"
	"github.com/uopensail/recgo-engine/pipeline/recalls"
	"github.com/uopensail/recgo-engine/userctx"
)

type IPipeline interface {
	GetName() string
	GetBuckets() []int
	Do(uCtx *userctx.UserContext) model.Collection
}

type Pipeline struct {
	name       string
	buckets    []int
	filter     freqs.IFilter
	recalls    []recalls.IRecall
	ranker     rank.IRank
	constrains constrains.IConstrains
}

func NewPipeline(conf *model.PipelineConfigure) *Pipeline {
	filter := freqs.NewFreqController(conf.Freqs)
	recallers := make([]recalls.IRecall, 0, len(conf.Recalls))
	for _, recall := range conf.Recalls {
		r := recalls.NewRecall(recall)
		if r != nil {
			recallers = append(recallers, r)
		}
	}
	ranker := rank.NewRank(conf.Rank)
	constrains := constrains.NewConstains(conf.Constrains)

	if filter == nil || len(recallers) == 0 || ranker == nil || constrains == nil {
		panic(fmt.Errorf("build pipeline fail"))
	}
	return &Pipeline{
		name:       conf.Name,
		buckets:    conf.Buckets,
		filter:     filter,
		recalls:    recallers,
		ranker:     ranker,
		constrains: constrains,
	}
}

func (p *Pipeline) GetName() string {
	return p.name
}

func (p *Pipeline) GetBuckets() []int {
	return p.buckets
}

func (p *Pipeline) Do(uCtx *userctx.UserContext) model.Collection {
	uCtx.Filter = p.filter.Do(uCtx)
	var wg sync.WaitGroup
	ch := make(chan model.Collection, len(p.recalls))
	for i := range p.recalls {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			collection := p.recalls[j].Do(uCtx)
			if collection != nil {
				ch <- collection
			}
		}(i)
	}

	wg.Wait()
	close(ch)

	maxSize := 0
	collections := make([]model.Collection, 0, len(p.recalls))
	for collection := range ch {
		collections = append(collections, collection)
		maxSize = max(maxSize, len(collection))
	}

	filter := make(map[int]*model.Entry)
	recall := make([]*model.Entry, 0, maxSize)
	for i := range maxSize {
		for _, collection := range collections {
			if len(collection) <= i {
				continue
			}
			entry := collection[i]
			if first, ok := filter[entry.ID]; ok {
				first.MergeChans(entry)
				continue
			}
			filter[entry.ID] = entry
			recall = append(recall, entry)
		}
	}

	recall = recall[:min(512, len(recall))]
	rank := p.ranker.Do(uCtx, recall)
	return p.constrains.Do(uCtx, rank)
}
