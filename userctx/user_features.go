package userctx

import (
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/ulib/commonconfig"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"

	grpc "google.golang.org/grpc"
)

type UserFeatures struct {
	UFeat *sample.MutableFeatures
}

var fukuGrpcConn *grpc.ClientConn

func Init(rdConf commonconfig.RegisterDiscoveryConfig) {
	fukuGrpcConn, _ = utils.NewKratosGrpcConn(rdConf)
}

func converUserTFeature(uid string, apiReq *recapi.RecRequest) *sample.MutableFeatures {
	feat := sample.NewMutableFeatures()
	feat.Set("u_d_click_list", &sample.Strings{Value: []string{"item_id_4589", "item_id_4408"}})

	if apiReq != nil {
		for k, v := range apiReq.UserFeature {
			if v != nil {
				switch sample.DataType(v.Type) {
				case sample.Int64Type:
					if v.Value != nil {
						feat.Set(k, &sample.Int64{Value: v.Value.Iv})
					}
				case sample.Float32Type:
					if v.Value != nil {
						feat.Set(k, &sample.Float32{Value: v.Value.Fv})
					}
				case sample.StringType:
					if v.Value != nil {
						feat.Set(k, &sample.String{Value: v.Value.Sv})
					}
				case sample.Int64sType:
					if v.Value != nil {
						feat.Set(k, &sample.Int64s{Value: v.Value.Ivs})
					}
				case sample.Float32sType:
					if v.Value != nil {
						feat.Set(k, &sample.Float32s{Value: v.Value.Fvs})
					}
				case sample.StringsType:
					if v.Value != nil {
						feat.Set(k, &sample.Strings{Value: v.Value.Svs})
					}
				}
			}
		}
	}
	return feat
}
