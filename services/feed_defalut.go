package services

import (
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/uapi/sunmaoapi"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"

	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
)

func (srv *Services) feedDefaultRec(uCtx *userctx.UserContext, modelEntities *strategy.ModelEntities) (*sunmaoapi.RecResult, error) {

	abCaseV := uCtx.ABData.Get("feed.default.rec")
	ret := sunmaoapi.RecResult{
		UserId:   uCtx.ApiRequest.UserId,
		DeviceId: uCtx.ApiRequest.DeviceId,
		TraceId:  uCtx.ApiRequest.TraceId,
	}

	if abCaseV != nil {
		switch abCaseV.CaseValue {
		case "base":
			uCtx.ABData.MarkHit(abCaseV.CaseId)
			ret.Expids = uCtx.Mark()
			//推荐数据埋点
			go reportLogsdk(uCtx, &ret)
			return &ret, nil
		case "exp":
			uCtx.ABData.MarkHit(abCaseV.CaseId)
			ret.Expids = uCtx.Mark()
		default:
			//推荐数据埋点
			go reportLogsdk(uCtx, &ret)
			return &ret, nil
		}
	}

	//do strategy
	istrategy := modelEntities.StrategyEntities.GetStrategy(int(HomeRecommendStrategyEntryID))

	if istrategy != nil {
		istrategy = strategy.BuildRuntimeEntity(modelEntities, uCtx, istrategy.Meta())
		recRes, err := istrategy.Do(uCtx)
		if err != nil {
			zlog.LOG.Error("strategy.do", zap.Error(err))
		} else {
			topk := int(uCtx.ApiRequest.Count)
			if len(recRes.StageList) < topk {
				prome.NewStat("rec_not_enough").SetCounter(topk - len(recRes.StageList)).End()
				topk = len(recRes.StageList)
			}
			items := recRes.StageList[:topk]
			ret.Items = make([]string, len(items))
			for i := 0; i < len(items); i++ {
				//TODO: 确定ID类型
				ret.Items[i] = items[i].Source.Key()
			}
		}
	}
	zlog.SLOG.Debug("api request response ", uCtx.ApiRequest, &ret)
	//推荐数据埋点
	go reportLogsdk(uCtx, &ret)
	return &ret, nil

}

func reportLogsdk(uCtx *userctx.UserContext, recRes *sunmaoapi.RecResult) {
	var itemStr buffer.Buffer

	for i := 0; i < len(recRes.Items); i++ {
		if i != 0 {
			itemStr.WriteString(",")
		}
		itemStr.WriteString(recRes.Items[i])
	}

	var expStr buffer.Buffer

	expStr.WriteString(recRes.Expids)

	//TODO:

	// logger.Push(&logger.Log{
	// 	ProductId: "honghu",
	// 	UserId:    uCtx.ApiRequest.UserId,
	// 	DeviceId:  uCtx.ApiRequest.DeviceId,
	// 	TraceId:   uCtx.ApiRequest.TraceId,
	// 	Ts:        time.Now().Unix(),
	// 	Eventid:   "rec_dist",
	// 	Items:     itemStr.String(),
	// 	Expids:    expStr.String(),
	// })
}
