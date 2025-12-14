package services

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/report"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Services contains service dependencies such as registry instance and report.
type Services struct {
	instance registry.ServiceInstance
	report   report.IReport
}

// NewServices creates a new Services instance and initializes report.
func NewServices() *Services {
	srv := Services{}
	srv.report = report.NewReport(config.AppConfigInstance.ReportConfig)
	return &srv
}

// RegisterGinRouter registers HTTP routes for the recommendation API.
func (srv *Services) RegisterGinRouter(ginEngine *gin.Engine) {
	apiV1 := ginEngine.Group("api/v1")
	{
		apiV1.POST("/feeds", srv.FeedsHandler)
		apiV1.POST("/related", srv.RelatedHandler)
	}
}

// FeedsHandler handles feed recommendations requests.
func (srv *Services) FeedsHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("HTTP.FeedsHandler")
	defer pStat.End()

	ctx, cancel := context.WithTimeout(gCtx.Request.Context(), 100*time.Millisecond)
	defer cancel()

	var req recapi.Request
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, recapi.Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	uCtx := userctx.NewUserContext(ctx, &req)
	resp := strategy.StrategyInstance.Feeds(uCtx)

	if resp == nil {
		gCtx.JSON(http.StatusInternalServerError, recapi.Response{
			Code:    -1,
			Message: "internal error: nil response",
		})
		return
	}

	srv.report.Report(uCtx, resp)
	gCtx.JSON(http.StatusOK, resp)
}

// RelatedHandler handles related item recommendations requests.
func (srv *Services) RelatedHandler(gCtx *gin.Context) {
	pStat := prome.NewStat("HTTP.RelatedHandler")
	defer pStat.End()

	ctx, cancel := context.WithTimeout(gCtx.Request.Context(), 100*time.Millisecond)
	defer cancel()

	var req recapi.Request
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, recapi.Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	uCtx := userctx.NewUserContext(ctx, &req)
	resp := strategy.StrategyInstance.Related(uCtx)

	if resp == nil {
		gCtx.JSON(http.StatusInternalServerError, recapi.Response{
			Code:    -1,
			Message: "internal error: nil response",
		})
		return
	}

	srv.report.Report(uCtx, resp)
	gCtx.JSON(http.StatusOK, resp)
}

// --- gRPC Health Check Implementation ---

func (srv *Services) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (srv *Services) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (srv *Services) Close() {
	// cleanup resources if needed
}
