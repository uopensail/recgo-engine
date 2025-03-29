package report

import (
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
)

type IReport interface {
	Report(uCtx *userctx.UserContext, recRes *recapi.RecResult) error
	Close()
}

func NewReport(cfg config.ReportConfig) IReport {
	switch cfg.Type {
	case "segment":
		return NewSegmentReport(cfg.SegmentConfig)
	case "slslog":
		return NewSLSLogReport(cfg.SLSLogConfig)
	}
	return &ZLogReport{}
}
