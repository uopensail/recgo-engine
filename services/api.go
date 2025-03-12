package services

import (
	"context"
	"time"

	"net/http"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/uapi/sunmaoapi"
	"github.com/uopensail/ulib/prome"

	"github.com/gin-gonic/gin"
)

func (srv *Services) Recommend(ctx context.Context, in *recapi.RecRequest) (*recapi.RecResponse, error) {
	stat := prome.NewStat("HomeRecommend")
	defer stat.End()
	modelEntities := strategy.EntitiesMgr.GetModelEntities()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	uCtx := userctx.NewUserContext(ctx, in, &modelEntities.Ress, &modelEntities.Model,
		&modelEntities.FilterResources)
	recResult, err := srv.feedDefaultRec(uCtx, modelEntities)
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
	pStat := prome.NewStat("RecommendHandler")
	defer pStat.End()

	var postData sunmaoapi.RecRequest
	if err := gCtx.ShouldBind(&postData); err != nil {
		gCtx.JSON(http.StatusInternalServerError, model.StatusResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	ret, err := srv.HomeRecommend(context.Background(), &postData)
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

	var postData recapi.RecRequest
	if err := gCtx.ShouldBind(&postData); err != nil {
		gCtx.JSON(http.StatusInternalServerError, model.StatusResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	modelEntities := strategy.EntitiesMgr.GetModelEntities()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	uCtx := userctx.NewUserContext(ctx, &postData, &modelEntities.Ress, &modelEntities.Model,
		&modelEntities.FilterResources)

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
