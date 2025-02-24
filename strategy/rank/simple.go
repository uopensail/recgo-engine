package rank

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/userctx"
)

func init() {
	RegisterPlugin(SimpleRankStrategy, NewSimpleRankEntity)
}

type SimpleRankEntity struct {
	cfg table.RankEntityMeta
}

func NewSimpleRankEntity(cfg table.RankEntityMeta, env config.EnvConfig) IStrategyEntity {
	return &SimpleRankEntity{cfg: cfg}
}
func (entity *SimpleRankEntity) Meta() *table.RankEntityMeta {
	return &entity.cfg
}
func (entity *SimpleRankEntity) Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error) {
	return in, nil
}
