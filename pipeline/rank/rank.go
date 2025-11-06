package rank

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
)

type IRank interface {
	Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection
}

func NewRank(conf model.IRank) IRank {
	switch conf.GetType() {
	case model.RankTypeChannelPriority:
		return NewChannelPriority(conf.(*model.ChannelPriorityRankConfigure))
	case model.RankTypeRule:
		return NewRule(conf.(*model.RuleBasedRankConfigure))
	case model.RankTypeModel:
		return NewModeler(conf.(*model.ModelBasedRankConfigure))
	}
	return nil
}
