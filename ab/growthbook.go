package ab

import (
	"context"
	"log"

	gb "github.com/growthbook/growthbook-golang"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type GrowthBookSDK struct {
	client *gb.Client
}

func newGrowthBookSDK(sdkConfig config.GrowthBookSDKConfig) *GrowthBookSDK {
	var opts []gb.ClientOption
	opts = append(opts, gb.WithSseDataSource())
	opts = append(opts, gb.WithClientKey(sdkConfig.ClientKey))

	if len(sdkConfig.APIHost) > 0 {
		opts = append(opts, gb.WithApiHost(sdkConfig.APIHost))
	}
	client, err := gb.NewClient(
		context.Background(),
		opts...,
	)
	if err != nil {
		zlog.LOG.Error("[AB.growthbook] NewClient error",
			zap.Error(err))
		return &GrowthBookSDK{}
	}
	if err := client.EnsureLoaded(context.Background()); err != nil {
		zlog.LOG.Fatal("Data loading failed: ", zap.Error(err))
	}
	return &GrowthBookSDK{client}
}
func (gbSDK *GrowthBookSDK) RequestABInfo(ctx context.Context, id string, feature *sample.MutableFeatures) ABInfo {
	return newGrowthBookABInfo(gbSDK.client, id, feature)
}
func (gbSDK *GrowthBookSDK) Close() {
	if gbSDK.client != nil {
		gbSDK.client.Close()
	}
}

type GrowthBookABInfo struct {
	childCli *gb.Client
}

func (ab *GrowthBookABInfo) EvalFeatureValue(ctx context.Context, featureKey string) string {
	if ab.childCli != nil {
		return ""
	}
	featureValue := ab.childCli.EvalFeature(ctx, featureKey)
	return featureValue.Value.(string)
}
func newGrowthBookABInfo(client *gb.Client, id string, feature *sample.MutableFeatures) *GrowthBookABInfo {
	if client == nil {
		return &GrowthBookABInfo{}
	}
	//TODO: feature toattrs
	attrs := gb.Attributes{"id": id}

	child, err := client.WithAttributes(attrs)
	if err != nil {
		log.Fatal("Child client creation failed: ", err)
	}
	return &GrowthBookABInfo{child}
}
