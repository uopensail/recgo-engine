package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"testing"
	"time"

	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/ulib/sample"
	"google.golang.org/grpc"
)

func Test_main(t *testing.T) {
	run("./conf/local/config.toml", "./logs")
	time.Sleep(time.Second * 3)

	conn, err := grpc.NewClient("localhost:3527", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	cli := recapi.NewRecServiceClient(conn)
	time.Sleep(time.Second * 3)
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
