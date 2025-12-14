package rank

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IRank defines the interface for ranking strategies.
// It takes a user context and a collection of entries, and returns a reordered collection.
type IRank interface {
	Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection
}

// NewRank creates a specific IRank implementation based on the given configuration type.
// Supported types:
// - model.RankTypeChannelPriority -> ChannelPriority ranker
// - model.RankTypeRule            -> Rule-based ranker
// - model.RankTypeModel           -> Model-based ranker
func NewRank(conf model.IRank) IRank {
	switch conf.GetType() {
	case model.RankTypeChannelPriority:
		if c, ok := conf.(*model.ChannelPriorityRankConfigure); ok {
			return NewChannelPriority(c)
		}
		zlog.LOG.Error("NewRank.TypeAssertError", zap.String("expected", "ChannelPriorityRankConfigure"))
	case model.RankTypeRule:
		if c, ok := conf.(*model.RuleBasedRankConfigure); ok {
			return NewRule(c)
		}
		zlog.LOG.Error("NewRank.TypeAssertError", zap.String("expected", "RuleBasedRankConfigure"))
	case model.RankTypeModel:
		if c, ok := conf.(*model.ModelBasedRankConfigure); ok {
			return NewModeler(c)
		}
		zlog.LOG.Error("NewRank.TypeAssertError", zap.String("expected", "ModelBasedRankConfigure"))
	default:
		zlog.LOG.Warn("NewRank.UnknownType", zap.String("type", conf.GetType()))
	}
	return nil
}
