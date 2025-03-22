package services

import (
	"context"

	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"

	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
)

func (srv *Services) feedDefaultRec(ctx context.Context, in *recapi.RecRequest, entities *strategy.ModelEntities) (*recapi.RecResult, error) {
	abInfo := userctx.FetchABInfo(userctx.UID(in))
	uCtx := userctx.NewUserContext(ctx, in, abInfo, &entities.Ress, &entities.Model,
		&entities.FilterResources)
	abCaseV := uCtx.ABData.Get("feed.default.rec")
	ret := recapi.RecResult{
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
	istrategy := entities.StrategyEntities.GetStrategy(in.Pipeline)
	if istrategy != nil {
		uCtx.InitUserFilter(istrategy.Meta().SubPoolID)

		istrategy = strategy.BuildRuntimeEntity(entities, uCtx, istrategy.Meta())
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

func reportLogsdk(uCtx *userctx.UserContext, recRes *recapi.RecResult) {
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
