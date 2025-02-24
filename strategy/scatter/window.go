package scatter

import (
	"encoding/json"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"go.uber.org/zap"

	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
)

func init() {
	RegisterPlugin("window", NewWindowEntity)
}

type GroupLimitConfig struct {
	Field string `json:"field"`
	Limit int    `json:"limit"`
}

type WindowConfig struct {
	GroupLimits []GroupLimitConfig `json:"group_limit"`
	WindowSize  int                `json:"window_size"`
}

type StageListState struct {
	OriginStageList model.RecItemList
	CurStageList    model.RecItemList
	Insert          map[int]int
}

func (stateList *StageListState) Index(i int) *model.ItemRefScore {
	if i < 0 {
		return nil
	}
	if i < len(stateList.CurStageList) {
		return &stateList.CurStageList[i]
	} else {
		j := len(stateList.CurStageList) - len(stateList.Insert)

		for k := 0; j < len(stateList.OriginStageList); j++ {
			if _, ok := stateList.Insert[j]; ok {
				continue
			}
			if k == i-len(stateList.CurStageList) {
				return &stateList.OriginStageList[j]
			}
			k++
		}
	}
	return nil
}

func (stateList *StageListState) Next(begin int) (*model.ItemRefScore, int) {
	for j := begin; j < len(stateList.OriginStageList); j++ {
		if _, ok := stateList.Insert[j]; ok {
			continue
		}
		return &stateList.OriginStageList[j], j
	}
	return nil, begin
}

type WindowState struct {
	groupLimitState map[string]map[string]int

	cacheStringFromat map[int]map[string][]string
	StageListState
}

func (wstate *WindowState) checkLimit(groupLimitCfg []GroupLimitConfig, item *model.ItemRefScore) bool {

	stringFromat := wstate.cacheFieldValueString(item)

	for _, glimit := range groupLimitCfg {
		gField := glimit.Field
		limit := glimit.Limit
		fieldValues := stringFromat[gField]
		curGFieldStat := wstate.groupLimitState[gField]
		for i := 0; i < len(fieldValues); i++ {
			value := fieldValues[i]
			gFieldV, _ := curGFieldStat[value]
			if gFieldV+1 > limit {
				return true
			}
		}

	}
	return false
}

func (wstate *WindowState) getStringFromatCache(itemIndex int, field string) []string {
	if v, ok := wstate.cacheStringFromat[itemIndex]; ok {
		if v2, ok2 := v[field]; ok2 {
			return v2
		}
	}
	return nil
}

func NewWindowState(gLimitCfgs []GroupLimitConfig) *WindowState {
	state := WindowState{
		groupLimitState:   make(map[string]map[string]int, len(gLimitCfgs)),
		cacheStringFromat: make(map[int]map[string][]string),
	}
	for _, gLimit := range gLimitCfgs {
		state.groupLimitState[gLimit.Field] = make(map[string]int)
	}
	return &state
}
func (wstate *WindowState) cacheFieldValueString(item *model.ItemRefScore) map[string][]string {
	if _, ok := wstate.cacheStringFromat[item.ID]; !ok {
		wstate.cacheStringFromat[item.ID] = make(map[string][]string)
		for k := range wstate.groupLimitState {
			wstate.cacheStringFromat[item.ID][k] = featureToStrings(item, k)
		}
	}

	return wstate.cacheStringFromat[item.ID]
}
func featureToStrings(item *model.ItemRefScore, field string) []string {
	feat := item.Get(field)
	if feat == nil {
		return nil
	}
	switch feat.Type() {
	case sample.Float32Type, sample.Float32sType:
		return nil
	case sample.StringType:
		v, _ := feat.GetString()
		return []string{v}
	case sample.StringsType:
		v, _ := feat.GetStrings()
		return v
	case sample.Int64Type:
		v, _ := feat.GetInt64()
		return []string{utils.Int642String(v)}
	case sample.Int64sType:
		vv, _ := feat.GetInt64s()
		ret := make([]string, len(vv))
		for i := 0; i < len(ret); i++ {
			ret[i] = utils.Int642String(vv[i])
		}
		return ret
	}
	return nil
}

func (wstate *WindowState) do(item *model.ItemRefScore) {
	for k, v := range wstate.groupLimitState {
		vv := wstate.getStringFromatCache(item.ID, k)
		if len(vv) > 0 {
			for i := 0; i < len(vv); i++ {
				v[vv[i]]++
			}
		}
	}
}

func (wstate *WindowState) undo(item *model.ItemRefScore) {
	for k, v := range wstate.groupLimitState {
		vv := wstate.getStringFromatCache(item.ID, k)
		if len(vv) > 0 {
			for i := 0; i < len(vv); i++ {
				v[vv[i]]--
			}
		}
	}
}

type WindowEntity struct {
	cfg table.ScatterEntityMeta

	windowCfg WindowConfig
}

func NewWindowEntity(cfg table.ScatterEntityMeta, env config.EnvConfig) IStrategyEntity {
	entity := &WindowEntity{
		cfg: cfg,
	}
	err := json.Unmarshal(cfg.PluginParams, &entity.windowCfg)
	if err != nil {
		zlog.LOG.Error("Parse PluginParams", zap.Error(err))
		return nil
	}
	return entity

}

func (entity *WindowEntity) Meta() *table.ScatterEntityMeta {
	return &entity.cfg
}

func (entity *WindowEntity) slideWindow(from int, wstate *WindowState) {

	item, begin := wstate.Next(from)
	for item != nil {
		// try add item
		if wstate.checkLimit(entity.windowCfg.GroupLimits, item) == false {
			//add success
			wstate.do(item)
			zlog.LOG.Debug("scatter.slider", zap.Int("from", begin), zap.Int("to", len(wstate.StageListState.CurStageList)), zap.Any("item", item.ItemFeatures.Source))
			if begin != from {
				wstate.StageListState.Insert[begin] = from
			}
			wstate.StageListState.CurStageList = append(wstate.StageListState.CurStageList, *item)

			break
		}
		next := begin + 1
		item, begin = wstate.Next(next)
	}

}

func (entity *WindowEntity) Do(uCtx *userctx.UserContext, in model.StageResult) (model.StageResult, error) {
	stat := prome.NewStat("scatter.WindowEntity.Do")
	defer stat.End()
	wstate := NewWindowState(entity.windowCfg.GroupLimits)
	wstate.StageListState.OriginStageList = in.StageList
	wstate.StageListState.Insert = make(map[int]int)
	from := -entity.windowCfg.WindowSize
	for i := 0; i < int(uCtx.ApiRequest.Count)*5 && i < len(in.StageList); i++ {
		beginCheck := len(wstate.CurStageList) - len(wstate.Insert)
		entity.slideWindow(beginCheck, wstate)
		undoItem := wstate.Index(from)
		zlog.LOG.Debug("scatter.slider", zap.Int("to", from))
		if undoItem != nil {
			wstate.undo(undoItem)
		}
		from++
	}

	//merge
	begin := len(wstate.StageListState.CurStageList) - len(wstate.StageListState.Insert)
	for i := begin; i < len(wstate.OriginStageList); i++ {
		if _, ok := wstate.StageListState.Insert[i]; ok {
			continue
		}
		wstate.StageListState.CurStageList = append(wstate.StageListState.CurStageList, wstate.OriginStageList[i])
	}
	in.StageList = wstate.StageListState.CurStageList
	return in, nil
}
