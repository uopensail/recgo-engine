package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"testing"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/uapi/sunmaoapi"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
)

func Test_main(t *testing.T) {
	run("./conf/local/config.toml", "./logs")

	conn, _ := utils.NewKratosGrpcConn(config.AppConfigInstance.ServerConfig.RegisterDiscoveryConfig)
	cli := sunmaoapi.NewRecServiceClient(conn)
	fmt.Println(cli.HomeRecommend(context.Background(), &sunmaoapi.RecRequest{
		UserId: "",
		Count:  10,
		UserFeature: map[string]*sunmaoapi.Feature{
			"u_d_click_list": {
				Type: int32(sample.StringsType),
				Value: &sunmaoapi.FeatureValue{
					Svs: []string{"item_id_1599", "item_id_4589", "item_id_4408"},
				},
			},
			"u_s_country": {
				Type: int32(sample.StringType),
				Value: &sunmaoapi.FeatureValue{
					Sv: "ctryus",
				},
			},
			"u_s_language": {
				Type: int32(sample.StringType),
				Value: &sunmaoapi.FeatureValue{
					Sv: "langen",
				},
			},
		},
	}))
	select {}
}
