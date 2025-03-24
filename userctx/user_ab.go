package userctx

import (
	"context"
	"strconv"

	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"

	gb "github.com/growthbook/growthbook-golang"
	"github.com/uopensail/kongming-sdk-go/sdkcore"
)

type IAB interface {
	EvalFeatureValue(featureKey string) string
}

type KongmingAB struct {
	sdkcore.ABData
}

func fetchKongmingABInfo(id string) *sdkcore.ABData {

	abInfo := sdkcore.ABSDK.Get(id, nil)
	if abInfo != nil {
		return abInfo
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}

func NewKongmingAB(id string) *KongmingAB {
	abInfo := fetchKongmingABInfo(id)
	if abInfo != nil {
		return &KongmingAB{*abInfo}
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}
func (ab *KongmingAB) EvalFeatureValue(featureKey string) string {
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

var growthBookSDK GrowthBookSDK

type GrowthBookSDK struct {
	client *gb.Client
}

type GrowthBookAB struct {
}

func NewGrowthBookAB() *GrowthBookAB {
	client, err := gb.NewClient(
		context.Background(),
		gb.WithClientKey("sdk-XXXX"),
		gb.WithSseDataSource(),
	)
	defer client.Close()
}

type UserAB struct {
	IAB
}

func NewUserAB(id string, feature *sample.MutableFeatures) UserAB {
	userAB := UserAB{}

	return userAB
}
