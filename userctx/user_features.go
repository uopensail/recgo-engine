package userctx

import (
	"github.com/uopensail/uapi/sunmaoapi"
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

func (uCtx *UserFeatures) initUserTFeature(uid string, apiReq *sunmaoapi.RecRequest) {
	feat := sample.NewMutableFeatures()
	feat.Set("u_d_click_list", &sample.Strings{Value: []string{"item_id_4589", "item_id_4408"}})
	uCtx.UFeat = feat
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

	// if nil == fukuGrpcConn {
	// 	return
	// }
	// realcli := fukuapi.NewFuKuClient(fukuGrpcConn)

	// req := fukuapi.FuKuRequest{
	// 	UserID:   uid,
	// 	Group:    "all",
	// 	Features: nil,
	// }

	// uFeat, err := realcli.Get(context.Background(), &req)
	// //fmt.Println(uFeat, err)
	// // uFeat, err := core.FukuSDK.Get(uid, "", nil)
	// if err != nil {
	// 	zlog.LOG.Error("FukuSDK.Get", zap.Error(err))
	// 	return
	// }
	// zlog.SLOG.Debug(uid, uFeat)

}
