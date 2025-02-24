package userctx

import (
	"context"

	"github.com/uopensail/recgo-engine/model/dbmodel"
	fresource "github.com/uopensail/recgo-engine/strategy/filter/resource"
	"github.com/uopensail/recgo-engine/strategy/recalls/resource"
	"github.com/uopensail/uapi/sunmaoapi"
	"github.com/uopensail/ulib/pool"
)

type UserContext struct {
	context.Context
	*dbmodel.DBTabelModel // 配置引用
	*pool.Pool
	FilterRess *fresource.Resources
	RecallRess *resource.Resources

	UserFeatures
	UserAB

	ApiRequest  *sunmaoapi.RecRequest
	ExcludeList []int
	RelateItem  *pool.Features
}

func NewUserContext(ctx context.Context, apiReq *sunmaoapi.RecRequest, dbModel *dbmodel.DBTabelModel,
	pl *pool.Pool, fress *fresource.Resources,
	rress *resource.Resources) *UserContext {
	uCtx := UserContext{
		Context:      ctx,
		ApiRequest:   apiReq,
		DBTabelModel: dbModel,
		Pool:         pl,
		FilterRess:   fress,
		RecallRess:   rress,
	}

	//tran
	if uCtx.ApiRequest != nil {
		for i := 0; i < len(uCtx.ApiRequest.ExcludeItems); i++ {
			item := pl.GetByKey(uCtx.ApiRequest.ExcludeItems[i])
			uCtx.ExcludeList = append(uCtx.ExcludeList, item.ID)
		}
		uCtx.RelateItem = pl.GetByKey(uCtx.ApiRequest.RelateItem)
	}

	//初始化ab信息
	uCtx.UserAB.initAB(uCtx.UID())
	//初始化用户特征
	uCtx.UserFeatures.initUserTFeature(uCtx.UID(), apiReq)

	return &uCtx
}

func (uCtx *UserContext) UID() string {
	if len(uCtx.ApiRequest.UserId) > 0 {
		return uCtx.ApiRequest.UserId
	}
	return uCtx.ApiRequest.DeviceId
}
