package rank

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

// Rule applies a rule-based scoring to reorder collection entries.
// The rule is evaluated via a Minia expression using entry runtime features and user context.
type Rule struct {
	conf    *model.RuleBasedRankConfigure // rank configuration
	program *minia.Minia                  // compiled Minia rule program
}

// NewRule creates a new Rule ranker with the provided configuration.
// It compiles the rule string into a Minia program: "result=<rule>".
func NewRule(conf *model.RuleBasedRankConfigure) *Rule {
	pStat := prome.NewStat("Rank.NewRule")
	defer pStat.End()
	if conf.Rule == "" {
		zlog.LOG.Warn("NewRule.EmptyRule", zap.String("name", conf.Name))
	}
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Rule)})
	return &Rule{
		conf:    conf,
		program: program,
	}
}

// Do executes the rule-based scoring process:
// 1. Evaluate the configured rule for each entry using its runtime features and user context.
// 2. Set the entry's KeyScore.Score.
// 3. Sort the collection by score in descending order, preserving relative order of equal scores.
func (rule *Rule) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("Rank.Rule.Do")
	defer pStat.End()

	for _, entry := range collection {
		// Evaluate rule
		result := rule.program.Eval(entry.Runtime.Basic, uCtx.Features, entry.Runtime.RunTime)
		scoreFeature := result.Get("result")

		// Default score
		entry.KeyScore.Score = 0

		// Extract score from feature
		if scoreFeature != nil {
			val, err := scoreFeature.GetFloat32()
			if err == nil {
				entry.KeyScore.Score = val
			} else {
				zlog.LOG.Warn("Rule.Do.ScoreExtractError", zap.String("key", entry.KeyScore.Key), zap.Error(err))
			}
		}

		// Debug log
		zlog.LOG.Debug("Rule.Do.EntryScore",
			zap.String("key", entry.KeyScore.Key),
			zap.Float32("score", entry.KeyScore.Score))
	}

	// Sort by score (descending), stable to preserve original order for equal scores
	sort.Stable(collection)

	zlog.LOG.Debug("Rule.Do.Completed",
		zap.Int("total", len(collection)),
		zap.String("rule", rule.conf.Rule))
	return collection
}
