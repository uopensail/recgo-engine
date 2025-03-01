package userctx

import (
	"context"

	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/resources"
	fresource "github.com/uopensail/recgo-engine/strategy/freqfilter/resource"
	"github.com/uopensail/uapi/sunmaoapi"
	"github.com/uopensail/ulib/pool"
)

type SubPool struct {
	collection resources.Collection
}
type WhiteList struct {
	Set map[int]struct{}
}

func NewWhitList(cs ...resources.Collection) *WhiteList {
	wh := WhiteList{}
	for _, c := range cs {
		for _, id := range c {
			wh.Set[id] = struct{}{}
		}
	}
	return &wh
}

type UserFilter struct {
	excludeCollection resources.Collection //黑明单

	whiteList *WhiteList
}

func newUserFilter(excludeList resources.Collection, subPool *SubPool, pl *pool.Pool, condition string) *UserFilter {
	userFilter := UserFilter{
		excludeCollection: excludeList,
	}
	var apiFilterStaticCollection resources.Collection
	if subPool != nil {
		if len(condition) != 0 {
			apiFilterStaticCollection = resources.BuildCollection(pl, subPool.collection, "", condition)
			userFilter.whiteList = NewWhitList(apiFilterStaticCollection)
		} else {
			userFilter.whiteList = NewWhitList(subPool.collection)
		}
	} else {
		if len(condition) != 0 {
			apiFilterStaticCollection = resources.BuildCollection(pl, pl.WholeCollection, "", condition)
			userFilter.whiteList = NewWhitList(apiFilterStaticCollection)
		}
		//不设置白名单,那么全部通过
	}

	return &userFilter
}

// true:表示合法,Pass false: 不合法
func (filter *UserFilter) Check(id int) bool {
	if resources.BinarySearch(filter.excludeCollection, id) {
		return false
	}
	//如果不设置，那么全部通过，比如没有subpool&&也没有conditon
	if filter.whiteList == nil {
		return true
	}
	//判断是否在子集中 如果在白名单子集中，那么代表合法
	if _, ok := filter.whiteList.Set[id]; !ok {
		return false
	}
	return true
}

type UserContext struct {
	context.Context

	*dbmodel.DBTabelModel // 配置引用
	Ress                  *resources.Resource
	FilterRess            *fresource.Resources

	UserFeatures
	UserAB
	UserFilter
	ApiRequest *sunmaoapi.RecRequest

	RelateItem *pool.Features
}

func NewUserContext(ctx context.Context, apiReq *sunmaoapi.RecRequest,
	ress *resources.Resource,
	dbModel *dbmodel.DBTabelModel,
	fress *fresource.Resources) *UserContext {
	uCtx := UserContext{
		Context:      ctx,
		ApiRequest:   apiReq,
		DBTabelModel: dbModel,
		Ress:         ress,
		FilterRess:   fress,
	}

	//tran
	var excludeList resources.Collection
	if uCtx.ApiRequest != nil {
		excludeList = make(resources.Collection, len(uCtx.ApiRequest.ExcludeItems))
		for i := 0; i < len(uCtx.ApiRequest.ExcludeItems); i++ {
			item := ress.Pool.GetByKey(uCtx.ApiRequest.ExcludeItems[i])
			excludeList = append(excludeList, item.ID)
		}
		uCtx.RelateItem = ress.Pool.GetByKey(uCtx.ApiRequest.RelateItem)
	}

	//初始化ab信息
	uCtx.UserAB.initAB(uCtx.UID())
	//初始化用户特征
	uCtx.UserFeatures.initUserTFeature(uCtx.UID(), apiReq)

	//查看命中的物料子集
	//TODO: 根据Pipeline 配置子集, 根据api conditon
	uCtx.UserFilter = *newUserFilter(excludeList, nil, ress.Pool, "")
	return &uCtx
}

func (uCtx *UserContext) UID() string {
	if len(uCtx.ApiRequest.UserId) > 0 {
		return uCtx.ApiRequest.UserId
	}
	return uCtx.ApiRequest.DeviceId
}
