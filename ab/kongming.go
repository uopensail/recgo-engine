package ab

import (
	"context"
	"strconv"

	"github.com/uopensail/kongming-sdk-go/sdkcore"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type KongmingABInfo struct {
	sdkcore.ABData
}

func fetchKongmingABInfo(id string, attr map[string]string) *sdkcore.ABData {

	abInfo := sdkcore.ABSDK.Get(id, attr)
	if abInfo != nil {
		return abInfo
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}

func newKongmingABInfo(id string, attr map[string]string) *KongmingABInfo {
	abInfo := fetchKongmingABInfo(id, attr)
	if abInfo != nil {
		return &KongmingABInfo{*abInfo}
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}
func (ab *KongmingABInfo) EvalFeatureValue(ctx context.Context, featureKey string) string {
	if ab != nil {
		layerID, err := strconv.Atoi(featureKey)
		if err != nil {
			zlog.LOG.Error("[AB.kongming] EvalFeatureValue input featureKey error",
				zap.String("featureKey", featureKey))
			return ""
		}
		expInfo := ab.ABData.GetByLayerID(layerID)
		if expInfo != nil {
			return expInfo.CaseValue
		}
	}
	return ""
}

type KongMingSDK struct {
	core sdkcore.KongMingABSDK
}

func newKongMingSDK(cfg sdkcore.KongMingSDKConfig) *KongMingSDK {
	sdk := &KongMingSDK{core: sdkcore.KongMingABSDK{}}
	sdk.core.Init(cfg)
	return sdk
}
func (sdk *KongMingSDK) RequestABInfo(ctx context.Context, id string, feature *sample.MutableFeatures) ABInfo {
	//TODO: feature to map[string]string
	return newKongmingABInfo(id, nil)
}
func (sdk *KongMingSDK) Close() {
}
