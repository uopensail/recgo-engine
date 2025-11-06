package report

import (
	"encoding/json"

	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type ZLogReport struct {
}

func (report *ZLogReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	zlog.LOG.Info("rec_dist", zap.String("data", string(data)))
	return nil
}
func (report *ZLogReport) Close() {

}
