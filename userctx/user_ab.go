package userctx

import (
	"github.com/uopensail/ulib/zlog"

	"github.com/uopensail/kongming-sdk-go/sdkcore"
)

type UserAB struct {
	sdkcore.ABData
}

func FetchABInfo(id string) *sdkcore.ABData {
	//TODO: Support  GrowthBook
	abInfo := sdkcore.ABSDK.Get(id, nil)
	if abInfo != nil {
		return abInfo
	} else {
		zlog.LOG.Error("GetABData nil")
		return nil
	}
}
