package report

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/analytics-go"
	"go.uber.org/zap"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
)

// SegmentReport implements IReport using Segment.io's analytics client.
// It sends recommendation results as "rec_dist" events to the configured Segment workspace.
type SegmentReport struct {
	cli analytics.Client
}

// NewSegmentReport creates a new SegmentReport using the provided Segment configuration.
// It returns nil if the client could not be created.
//
// Example:
//
//	cfg := config.SegmentConfig{WriteKey: "xxx", Endpoint: "https://..."}
//	report := report.NewSegmentReport(cfg)
//	defer report.Close()
func NewSegmentReport(cfg config.SegmentConfig) *SegmentReport {
	client, err := analytics.NewWithConfig(cfg.WriteKey, analytics.Config{
		Endpoint: cfg.Endpoint,
	})
	if err != nil {
		zlog.LOG.Warn("failed to create Segment client", zap.Error(err))
		return nil
	}
	return &SegmentReport{cli: client}
}

// Report sends the recommendation response to Segment.io for analytics tracking.
// The data is marshaled to JSON and sent as a property in the "rec_dist" event.
func (report *SegmentReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	if report.cli == nil {
		return fmt.Errorf("Segment client is not initialized")
	}

	// Validate input
	if uCtx == nil || uCtx.Request.UserId == "" {
		return fmt.Errorf("invalid user context")
	}
	if resp == nil {
		return fmt.Errorf("nil recommendation response")
	}

	// Serialize response to JSON
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Prepare event properties
	properties := analytics.NewProperties().Set("data", string(data))

	// Send event to Segment
	return report.cli.Enqueue(analytics.Track{
		UserId:     uCtx.Request.UserId,
		Event:      "rec_dist",
		Properties: properties,
	})
}

// Close releases the Segment client resources (flushes pending events).
func (report *SegmentReport) Close() {
	if report.cli != nil {
		report.cli.Close()
	}
}
