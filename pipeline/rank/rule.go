package rank

import (
	"errors"
	"sort"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/program"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Rule applies a rule-based scoring to reorder collection entries.
// The rule is evaluated via a Minia expression using entry runtime features and user context.
type Rule struct {
	conf    *model.RuleBasedRankConfigure // rank configuration
	program *program.Program              // compiled Minia rule program
}

// NewRule creates a new Rule ranker with the provided configuration.
// It compiles the rule string into a Minia program: "result=<rule>".
func NewRule(conf *model.RuleBasedRankConfigure) *Rule {
	pStat := prome.NewStat("Rank.NewRule")
	defer pStat.End()
	if conf.Rule == "" {
		msg := "Rank.NewRule.EmptyRule"
		zlog.LOG.Error(msg, zap.String("name", conf.Name))
		panic(errors.New(msg))
	}

	program, err := program.NewProgram(conf.Rule)
	if err != nil {
		zlog.LOG.Error("Rank.NewRule program create error",
			zap.Error(err))
		panic(err)
	}

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
		value, err := rule.program.Eval(entry.Runtime.Basic, uCtx.Features, entry.Runtime.RunTime)
		if err != nil {
			zlog.LOG.Error("Rank.Rule.Do program eval error",
				zap.Error(err))
		}

		// Default score
		entry.KeyScore.Score = 0.0
		score, ok := value.(float32)
		if ok {
			entry.KeyScore.Score = score
		} else {
			zlog.LOG.Error("Rule.Do.ScoreExtractError", zap.String("key", entry.KeyScore.Key))
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
