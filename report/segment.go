package report

import (
	"encoding/json"

	"github.com/segmentio/analytics-go"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
)

type SegmentReport struct {
	cli analytics.Client
}

func NewSegmentReport(cfg config.SegmentConfig) *SegmentReport {
	client, _ := analytics.NewWithConfig(cfg.WriteKey, analytics.Config{
		Endpoint: cfg.Endpoint,
	})

	return &SegmentReport{cli: client}
}

func (report *SegmentReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	if report.cli != nil {
		protpies := analytics.NewProperties()
		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		protpies.Set("data", data)

		report.cli.Enqueue(analytics.Track{
			UserId: uCtx.Request.UserId,
			Event:  "rec_dist",
		})
	}
	return nil
}
func (report *SegmentReport) Close() {
	if report.cli != nil {
		report.cli.Close()
	}
}
