package services

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/recgo-engine/utils"
	etcdclient "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type StrategyEntryID int

const (
	HomeRecommendStrategyEntryID   StrategyEntryID = 0
	RelateRecommendStrategyEntryID StrategyEntryID = 0
)

type Services struct {
	recapi.UnimplementedRecServiceServer
	etcdCli *etcdclient.Client

	instance registry.ServiceInstance
}

func NewServices() *Services {
	srv := Services{}

	return &srv
}
func (srv *Services) Init(configFolder string, etcdName string, etcdCli *etcdclient.Client, reg utils.Register) {
	srv.etcdCli = etcdCli
	jobUtil := utils.NewMetuxJobUtil(etcdName, reg, etcdCli, 10, -1)
	strategy.EntitiesMgr.Init(config.AppConfigInstance.EnvConfig, jobUtil)

}
func (srv *Services) RegisterGrpc(grpcS *grpc.Server) {
	recapi.RegisterRecServiceServer(grpcS, srv)
	//grpc_health_v1.RegisterHealthServer(grpcS, srv)
}

func (srv *Services) RegisterGinRouter(ginEngine *gin.Engine) {
	apiV1 := ginEngine.Group("api/v1")
	{
		apiV1.POST("/home_rec", srv.RecommendHandler)
	}
	ginEngine.POST("/user", srv.UsrCtxInfoHandler)
}

func (srv *Services) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (srv *Services) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (srv *Services) Close() {

}
