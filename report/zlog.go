package report

import (
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type ZLogReport struct {
}

func (report *ZLogReport) Report(uCtx *userctx.UserContext, recRes *recapi.RecResult) error {
	recReportMap := recRes.ToMap()

	zapField := make([]zap.Field, 0, len(recReportMap))
	for k, v := range recReportMap {
		zapField = append(zapField, zap.String(k, v))
	}
	zlog.LOG.Info("rec_dist", zapField...)
	return nil
}
func (report *ZLogReport) Close() {

}
