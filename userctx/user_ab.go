package userctx

import (
	"github.com/uopensail/ulib/zlog"

	"github.com/uopensail/kongming-sdk-go/sdkcore"
)

type UserAB struct {
	sdkcore.ABData
}

func (uCtx *UserAB) initAB(id string) {
	//TODO: Support  GrowthBook
	abInfo := sdkcore.ABSDK.Get(id, nil)
	if abInfo != nil {
		uCtx.ABData = *abInfo
	} else {
		zlog.LOG.Error("GetABData nil")
	}
}
