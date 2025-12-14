package report

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IReport defines the interface for reporting recommendation results.
// Different implementations can store or forward recommendation logs.
type IReport interface {
	// Report sends the recommendation result along with the user context.
	// Implementations should handle any serialization or transformation internally.
	Report(uCtx *userctx.UserContext, recRes *recapi.Response) error

	// Close releases any resources held by the reporter (e.g., flush buffers, close connections).
	Close()
}

// NewReport creates a new IReport implementation based on configuration.
// Supported types:
//   - "segment" : returns a SegmentReport
//   - "slslog"  : returns an SLSLogReport
//   - default   : returns a ZLogReport as a fallback
//
// Example:
//
//	cfg := config.ReportConfig{Type: "segment", SegmentConfig: ...}
//	reporter := report.NewReport(cfg)
//	defer reporter.Close()
//	reporter.Report(uCtx, recRes)
func NewReport(cfg config.ReportConfig) IReport {
	switch cfg.Type {
	case "segment":
		return NewSegmentReport(cfg.SegmentConfig)
	case "slslog":
		return NewSLSLogReport(cfg.SLSLogConfig)
	default:
		zlog.LOG.Warn("unsupported report type, using ZLogReport as default", zap.String("type", cfg.Type))
		return &ZLogReport{}
	}
}
