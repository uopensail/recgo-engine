package scatter

import (
	"encoding/json"
	"sort"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"

	"github.com/uopensail/recgo-engine/userctx"
)

func init() {
	RegisterPlugin("weight", NewWeightEntity)
}

type weightParams struct {
	PropertyWeight map[string]float32 `json:"weight"`
}
type WeightEntity struct {
	cfg table.ScatterEntityMeta
	weightParams
}

func NewWeightEntity(cfg table.ScatterEntityMeta, env config.EnvConfig) IStrategyEntity {
	entity := &WeightEntity{
		cfg: cfg,
	}
	json.Unmarshal([]byte(cfg.PluginParams), &entity.weightParams)
	return entity
}

func (entity *WeightEntity) Meta() *table.ScatterEntityMeta {
	return &entity.cfg
}

func (entity *WeightEntity) Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error) {
	scores := make([]float32, len(in.StageList))
	counter := make(map[string]map[string]int)
	for i := 0; i < len(in.StageList); i++ {
		item := in.StageList[i]
		score := float32(i)
		for p, w := range entity.weightParams.PropertyWeight {
			if _, ok := counter[p]; !ok {
				counter[p] = make(map[string]int)
			}
			feat := item.Source.Get(p)
			switch feat.Type() {
			case sample.Int64Type:

				if v, err := feat.GetInt64(); err == nil {
					score += float32(counter[p][utils.Int642String(v)]) * w
					counter[p][utils.Int642String(v)]++
				}

			case sample.StringType:
				if v, err := feat.GetString(); err == nil {
					score += float32(counter[p][v]) * w
					counter[p][v]++
				}
			case sample.Int64sType:
			case sample.StringsType:

			case sample.Float32Type:
				fallthrough
			case sample.Float32sType:
				fallthrough
			default:
				continue
			}
		}
		scores[i] = score
	}
	sort.SliceStable(in.StageList, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	return in, nil
}
