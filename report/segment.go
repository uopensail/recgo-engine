package report

import (
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

func (report *SegmentReport) Report(uCtx *userctx.UserContext, recRes *recapi.RecResult) error {
	if report.cli != nil {
		protpies := analytics.NewProperties()
		recReportMap := recRes.ToMap()
		for k, v := range recReportMap {
			protpies.Set(k, v)
		}
		report.cli.Enqueue(analytics.Track{
			UserId: recRes.UserId,
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
