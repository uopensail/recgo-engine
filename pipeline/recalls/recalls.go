package recalls

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
)

type IRecall interface {
	Do(uCtx *userctx.UserContext) model.Collection
}

func NewRecall(conf model.IRecall) IRecall {
	switch conf.GetType() {
	case model.RecallTypeMatch:
		return NewMatcher(conf.(*model.MatchRecallConfigure))
	case model.RecallTypeModel:
		return NewModeler(conf.(*model.ModelRecallConfigure))
	}
	return nil
}
