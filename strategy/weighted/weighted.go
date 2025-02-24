package weighted

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/userctx"
)

const (
	DefalutStrategy = "default"
)

type IStrategyEntity interface {
	Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error)
	Meta() *table.WeightedEntityMeta
}
