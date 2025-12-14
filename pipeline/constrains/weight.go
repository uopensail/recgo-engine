package constrains

import (
	"fmt"
	"sort"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// WeightAdjust modifies the score of entries based on a configurable condition.
// If an entry matches the condition, its score is multiplied by the specified ratio.
// This constraint is useful to promote or demote certain items while maintaining order.
type WeightAdjust struct {
	conf    *model.WeightAdjustedConstrainConfigure // configuration for weight adjustment
	program *minia.Minia                            // compiled condition program
}

// NewWeightAdjust constructs a WeightAdjust from the given configuration.
func NewWeightAdjust(conf *model.WeightAdjustedConstrainConfigure) *WeightAdjust {
	pStat := prome.NewStat("NewWeightAdjust")
	defer pStat.End()
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Condition)})
	return &WeightAdjust{
		conf:    conf,
		program: program,
	}
}

// Do evaluates each entry in the collection against the configured condition.
// If the condition evaluates to 1, the entry's score is adjusted by multiplying
// it with the configured ratio. After adjustment, the collection is sorted.
func (w *WeightAdjust) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("WeightAdjust.Do")
	defer pStat.End()
	for _, entry := range collection {
		// Make sure the parameter order matches other modules: basic, user features, runtime.
		value := w.program.Eval(entry.Runtime.Basic, uCtx.Features, entry.Runtime.RunTime)
		if value != nil {
			result := value.Get("result")
			if result == nil {
				continue
			}
			hit, err := result.GetInt64()
			if err == nil && hit == 1 {
				// Apply score adjustment
				oldScore := entry.KeyScore.Score
				entry.KeyScore.Score *= w.conf.Ratio

				zlog.LOG.Debug("WeightAdjust.Applied",
					zap.String("entry_key", entry.KeyScore.Key),
					zap.Float32("old_score", oldScore),
					zap.Float32("new_score", entry.KeyScore.Score),
					zap.Float32("ratio", w.conf.Ratio))
			}
		}
	}

	// Sort using stable sort to maintain order for equal scores
	sort.Stable(collection)

	zlog.LOG.Debug("WeightAdjust.Completed",
		zap.Int("total_entries", len(collection)),
		zap.String("condition", w.conf.Condition),
		zap.Float32("ratio", w.conf.Ratio))

	return collection
}
