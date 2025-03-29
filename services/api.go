package services

import (
	"context"
	"time"

	"net/http"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/prome"

	"github.com/gin-gonic/gin"
)

func (srv *Services) Recommend(ctx context.Context, in *recapi.RecRequest) (*recapi.RecResponse, error) {
	stat := prome.NewStat("GRPC.HomeRecommend")
	defer stat.End()
	entities := strategy.EntitiesMgr.GetEntities()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	inWrapper := recapi.RecRequestWrapper{
		RecRequest: in,
		FieldsType: entities.Ress.FieldDataType,
	}
	inWrapper.FromRecRequest(in)
	recResult, err := srv.feedDefaultRec(ctx, &inWrapper, &entities.ModelEntities)
	if err != nil {
		return nil, err
	}
	return &recapi.RecResponse{
		Code: 0,
		Msg:  "OK",
		Data: recResult,
	}, nil
}

// RecommendHandler @Summary 获取命中的实验
// @BasePath /api/v1
// @Accept  json
// @Produce  json
// @Param payload body sunmaoapi.RecRequest true "RecRequest"
// @Success 200 {object} sunmaoapi.RecResponse
// @Failure 500 {object} model.StatusResponse
// @Failure 400 {object} model.StatusResponse
// @Router /rec [post]
func (srv *Services) RecommendHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("HTTP.RecommendHandler")
	defer pStat.End()
	entities := strategy.EntitiesMgr.GetEntities()
	ctx, cancel := context.WithTimeout(gCtx.Request.Context(), time.Millisecond*100)
	defer cancel()

	var postData recapi.RecRequestWrapper
	if err := gCtx.ShouldBindJSON(&postData); err != nil {
		gCtx.JSON(http.StatusInternalServerError, model.StatusResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	ret, err := srv.feedDefaultRec(ctx, &postData, &entities.ModelEntities)

	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, model.StatusResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	gCtx.JSON(http.StatusOK, ret)
	return
}

// for debug
func (srv *Services) UsrCtxInfoHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("UsrCtxInfo")
	defer pStat.End()

	var postData recapi.RecRequestWrapper
	if err := gCtx.ShouldBind(&postData); err != nil {
		gCtx.JSON(http.StatusInternalServerError, model.StatusResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	entities := strategy.EntitiesMgr.GetEntities()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	uCtx := userctx.NewUserContext(ctx, &postData, &entities.Ress, 0, &entities.Model,
		&entities.FilterResources)

	uCtxInfo := struct {
		userctx.UserFeatures
		userctx.UserAB
	}{
		UserFeatures: uCtx.UserFeatures,
		UserAB:       uCtx.UserAB,
	}
	gCtx.JSON(http.StatusOK, uCtxInfo)
	return
}
