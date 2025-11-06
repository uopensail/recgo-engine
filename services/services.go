package services

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/report"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type Services struct {
	instance registry.ServiceInstance
	report   report.IReport
}

func NewServices() *Services {
	srv := Services{}
	srv.report = report.NewReport(config.AppConfigInstance.ReportConfig)
	return &srv
}

func (srv *Services) RegisterGinRouter(ginEngine *gin.Engine) {
	apiV1 := ginEngine.Group("api/v1")
	{
		apiV1.POST("/feeds", srv.FeedsHandler)
		apiV1.POST("/related", srv.RelatedHandler)
	}
}

func (srv *Services) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (srv *Services) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (srv *Services) Close() {
}
