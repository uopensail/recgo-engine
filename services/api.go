package services

import (
	"context"
	"time"

	"net/http"

	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/prome"

	"github.com/gin-gonic/gin"
)

func (srv *Services) FeedsHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("HTTP.FeedsHandler")
	defer pStat.End()
	ctx, cancel := context.WithTimeout(gCtx.Request.Context(), time.Millisecond*100)
	defer cancel()

	var req recapi.Request
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: -1,
			Msg:  err.Error(),
		})
	}

	uCtx := userctx.NewUserContext(ctx, &req)

	resp := strategy.StrategyInstance.Feeds(uCtx)

	if resp == nil {
		gCtx.JSON(http.StatusInternalServerError, struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: -1,
			Msg:  "error",
		})
		return
	}

	srv.report.Report(uCtx, resp)
	gCtx.JSON(http.StatusOK, resp)
	return
}

func (srv *Services) RelatedHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("HTTP.RelatedHandler")
	defer pStat.End()
	// entities := strategy.EntitiesMgr.GetEntities()
	ctx, cancel := context.WithTimeout(gCtx.Request.Context(), time.Millisecond*100)
	defer cancel()

	var req recapi.Request
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: -1,
			Msg:  err.Error(),
		})
	}

	uCtx := userctx.NewUserContext(ctx, &req)

	resp := strategy.StrategyInstance.Feeds(uCtx)

	if resp == nil {
		gCtx.JSON(http.StatusInternalServerError, struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: -1,
			Msg:  "error",
		})
		return
	}

	srv.report.Report(uCtx, resp)
	gCtx.JSON(http.StatusOK, resp)
	return
}
