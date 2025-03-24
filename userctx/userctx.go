package userctx

import (
	"context"

	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/resources"
	fresource "github.com/uopensail/recgo-engine/strategy/freqfilter/resource"
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

func newUserFilter(excludeList resources.Collection, subPool *SubPool, ress *resources.Resource, condition string) *UserFilter {
	userFilter := UserFilter{
		excludeCollection: excludeList,
	}
	var apiFilterStaticCollection resources.Collection
	if subPool != nil {
		if len(condition) != 0 {
			apiFilterStaticCollection = resources.BuildCollection(ress, subPool.collection, condition)
			userFilter.whiteList = NewWhitList(apiFilterStaticCollection)
		} else {
			userFilter.whiteList = NewWhitList(subPool.collection)
		}
	} else {
		if len(condition) != 0 {
			apiFilterStaticCollection = resources.BuildCollection(ress, ress.Pool.WholeCollection, condition)
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
	SubPool               *SubPool
	FilterRess            *fresource.Resources

	UserFeatures
	UserAB
	UserFilter
	ApiRequest *recapi.RecRequest

	RelateItem *pool.Features
}

func NewUserContext(ctx context.Context, apiReq *recapi.RecRequest,

	ress *resources.Resource,
	subPoolID int,
	dbModel *dbmodel.DBTabelModel,
	fress *fresource.Resources) *UserContext {

	uCtx := UserContext{
		Context:      ctx,
		ApiRequest:   apiReq,
		DBTabelModel: dbModel,
		Ress:         ress,
		FilterRess:   fress,
	}
	uFeat := converUserTFeature(uCtx.UID(), apiReq)
	uCtx.UserAB = NewUserAB(uCtx.UID(), uFeat)
	//初始化用户特征
	uCtx.UserFeatures.UFeat = uFeat
	//tran

	if uCtx.ApiRequest != nil {
		uCtx.RelateItem = ress.Pool.GetByKey(uCtx.ApiRequest.RelateItem)
	}

	var subPool *SubPool
	subCollection, ok := uCtx.Ress.SubPoolCollectionRess.SubPool[subPoolID]
	if ok {
		subPool = &SubPool{
			collection: subCollection,
		}
		uCtx.SubPool = subPool
	}
	uCtx.initUserFilter(uCtx.SubPool)
	return &uCtx
}

func (uCtx *UserContext) initUserFilter(subPool *SubPool) {
	var excludeList resources.Collection

	if uCtx.ApiRequest != nil {
		excludeList = make(resources.Collection, len(uCtx.ApiRequest.ExcludeItems))
		for i := 0; i < len(uCtx.ApiRequest.ExcludeItems); i++ {
			item := uCtx.Ress.Pool.GetByKey(uCtx.ApiRequest.ExcludeItems[i])
			excludeList = append(excludeList, item.ID)
		}
	}

	uCtx.UserFilter = *newUserFilter(excludeList, subPool, uCtx.Ress, uCtx.ApiRequest.FilterCondition)

}
func UID(apiReq *recapi.RecRequest) string {
	if len(apiReq.UserId) > 0 {
		return apiReq.UserId
	}
	return apiReq.DeviceId
}

func (uCtx *UserContext) UID() string {
	if len(uCtx.ApiRequest.UserId) > 0 {
		return uCtx.ApiRequest.UserId
	}
	return uCtx.ApiRequest.DeviceId
}
