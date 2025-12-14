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
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

const (
	// MaxRecallLimit defines the maximum number of items that can pass from recall to ranking stage.
	MaxRecallLimit = 512
)

// IPipeline defines the interface for a recommendation pipeline.
// A pipeline orchestrates filter, recall, rank, and constraints to produce final recommendations.
type IPipeline interface {
	GetName() string
	Do(uCtx *userctx.UserContext) model.Collection
}

// Pipeline implements the IPipeline interface, holding the configured stages:
// - Frequency filter
// - Multiple recall strategies (parallel)
// - Ranking strategy
// - Constraints processor
type Pipeline struct {
	name       string
	filter     freqs.IFilter
	recalls    []recalls.IRecall
	ranker     rank.IRank
	constrains constrains.IConstrains
}

// NewPipeline creates a new Pipeline from configuration.
// Panics if any mandatory stage is missing.
func NewPipeline(conf *model.PipelineConfigure) *Pipeline {
	pStat := prome.NewStat("NewPipeline")
	defer pStat.End()

	filter := freqs.NewFreqController(conf.Freqs)

	recallers := make([]recalls.IRecall, 0, len(conf.Recalls))
	for _, recallConf := range conf.Recalls {
		r := recalls.NewRecall(recallConf)
		if r != nil {
			recallers = append(recallers, r)
		}
	}

	ranker := rank.NewRank(conf.Rank)
	constrains := constrains.NewConstains(conf.Constrains)

	if filter == nil || len(recallers) == 0 || ranker == nil || constrains == nil {
		panic(fmt.Errorf("build pipeline fail: missing stage"))
	}

	zlog.LOG.Info("Pipeline.Created",
		zap.String("name", conf.Name),
		zap.Int("recallers_count", len(recallers)))

	return &Pipeline{
		name:       conf.Name,
		filter:     filter,
		recalls:    recallers,
		ranker:     ranker,
		constrains: constrains,
	}
}

// GetName returns the name of the pipeline.
func (p *Pipeline) GetName() string {
	return p.name
}

// Do executes the pipeline stages sequentially:
// 1. Filter stage
// 2. Parallel recall stage
// 3. Merge and deduplicate recalled items
// 4. Ranking stage
// 5. Constraints stage
func (p *Pipeline) Do(uCtx *userctx.UserContext) model.Collection {
	pStat := prome.NewStat("Pipeline.Do")
	defer pStat.End()

	// Step 1: Frequency filter
	uCtx.Filter = p.filter.Do(uCtx)

	// Step 2: Parallel recall
	var wg sync.WaitGroup
	ch := make(chan model.Collection, len(p.recalls))
	for i := range p.recalls {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			collection := p.recalls[idx].Do(uCtx)
			if collection != nil {
				ch <- collection
				zlog.LOG.Debug("Pipeline.RecallCompleted",
					zap.Int("recall_index", idx),
					zap.Int("items_count", len(collection)))
			}
		}(i)
	}
	wg.Wait()
	close(ch)

	// Step 3: Merge recalled collections with deduplication
	maxSize := 0
	collections := make([]model.Collection, 0, len(p.recalls))
	for col := range ch {
		collections = append(collections, col)
		if len(col) > maxSize {
			maxSize = len(col)
		}
	}

	filter := make(map[int]*model.Entry)
	recall := make([]*model.Entry, 0, maxSize)

	for i := 0; i < maxSize; i++ {
		for _, col := range collections {
			if len(col) <= i {
				continue
			}
			entry := col[i]
			if first, exists := filter[entry.ID]; exists {
				first.MergeChans(entry)
			} else {
				filter[entry.ID] = entry
				recall = append(recall, entry)
			}
		}
	}

	// Limit recall size to constant value
	if len(recall) > MaxRecallLimit {
		recall = recall[:MaxRecallLimit]
	}

	zlog.LOG.Debug("Pipeline.MergedRecall",
		zap.Int("merged_count", len(recall)),
		zap.Int("original_collections", len(collections)))

	// Step 4: Rank
	ranked := p.ranker.Do(uCtx, recall)
	zlog.LOG.Debug("Pipeline.Ranked", zap.Int("ranked_count", len(ranked)))

	// Step 5: Constraints
	final := p.constrains.Do(uCtx, ranked)
	zlog.LOG.Debug("Pipeline.Completed", zap.Int("final_count", len(final)))

	pStat.SetCounter(len(final))
	return final
}
