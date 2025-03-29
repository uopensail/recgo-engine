package ab

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/uopensail/kongming-sdk-go/sdkcore"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type KongmingABInfo struct {
	sdkcore.ABData
	hitInfo map[string]map[string]string
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
		return &KongmingABInfo{
			ABData:  *abInfo,
			hitInfo: map[string]map[string]string{},
		}
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}

func (ab *KongmingABInfo) HitInfo() string {
	data, err := json.Marshal(ab.hitInfo)
	if err != nil {
		return ""
	}
	return string(data)
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

			hitInfo := make(map[string]string)
			hitInfo["experimentId"] = utils.Int642String(expInfo.ExpId)
			hitInfo["variationId"] = utils.Int642String(expInfo.CaseId)
			hitInfo["featureKey"] = featureKey
			ab.hitInfo[featureKey] = hitInfo
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
