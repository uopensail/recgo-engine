package rank

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/userctx"
)

const (
	SimpleRankStrategy = "simple_rank"
)

type IStrategyEntity interface {
	Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error)
	Meta() *table.RankEntityMeta
}
