package weighted

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/recgo-engine/userctx"
)

func init() {
	RegisterPlugin(DefalutStrategy, NewDefaultStrategyEntity)
}

type DefaultStrategyEntity struct {
	cfg table.WeightedEntityMeta
}

func NewDefaultStrategyEntity(cfg table.WeightedEntityMeta, env config.EnvConfig) IStrategyEntity {
	return &DefaultStrategyEntity{
		cfg: cfg,
	}

}
func (entity *DefaultStrategyEntity) Meta() *table.WeightedEntityMeta {
	return &entity.cfg
}
func (entity *DefaultStrategyEntity) Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error) {
	return in, nil
}
