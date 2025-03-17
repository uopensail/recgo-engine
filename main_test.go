package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"testing"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
)

func Test_main(t *testing.T) {
	run("./conf/local/config.toml", "./logs")

	conn, _ := utils.NewKratosGrpcConn(config.AppConfigInstance.ServerConfig.RegisterDiscoveryConfig)
	cli := recapi.NewRecServiceClient(conn)
	fmt.Println(cli.Recommend(context.Background(), &recapi.RecRequest{
		UserId: "",
		Count:  10,
		UserFeature: map[string]*recapi.Feature{
			"u_d_click_list": {
				Type: int32(sample.StringsType),
				Value: &recapi.FeatureValue{
					Svs: []string{"item_id_1599", "item_id_4589", "item_id_4408"},
				},
			},
			"u_s_country": {
				Type: int32(sample.StringType),
				Value: &recapi.FeatureValue{
					Sv: "ctryus",
				},
			},
			"u_s_language": {
				Type: int32(sample.StringType),
				Value: &recapi.FeatureValue{
					Sv: "langen",
				},
			},
		},
	}))
	select {}
}
