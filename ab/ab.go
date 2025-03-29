package ab

import (
	"context"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/ulib/sample"
)

type ABInfo interface {
	EvalFeatureValue(ctx context.Context, featureKey string) string
	HitInfo() string
}

type ABCli interface {
	RequestABInfo(ctx context.Context, id string, feature *sample.MutableFeatures) ABInfo
	Close()
}

func InitABClient(abConfig config.ABConfig) ABCli {
	switch abConfig.Type {
	case "kongming":
		return newKongMingSDK(abConfig.KongMingSDKConfig)
	case "growthbook":
		return newGrowthBookSDK(abConfig.GrowthBookSDKConfig)
	}
	return &KongMingSDK{}
}

var ABClient ABCli
