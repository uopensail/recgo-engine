package userctx

import (
	"context"

	"github.com/uopensail/recgo-engine/ab"
	"github.com/uopensail/ulib/sample"
)

type UserAB struct {
	AbInfo ab.ABInfo
}

func NewUserAB(ctx context.Context, id string, feature *sample.MutableFeatures) UserAB {
	userAB := UserAB{}
	userAB.AbInfo = ab.ABClient.RequestABInfo(ctx, id, feature)
	return userAB
}
