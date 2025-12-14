package recalls

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IRecall defines the interface for recall strategies.
// The Do method retrieves a collection of candidate items based on the recall configuration.
type IRecall interface {
	Do(uCtx *userctx.UserContext) model.Collection
}

// NewRecall creates an IRecall implementation based on the provided configuration type.
// Supported types:
//   - model.RecallTypeMatch -> returns a Matcher-based recall
//   - model.RecallTypeModel -> returns a Modeler-based recall
//
// Returns nil if the type is unknown or if type assertion fails.
func NewRecall(conf model.IRecall) IRecall {
	recallType := conf.GetType()
	zlog.LOG.Info("NewRecall.Init", zap.String("type", recallType), zap.String("name", conf.GetName()))

	switch recallType {
	case model.RecallTypeMatch:
		if cfg, ok := conf.(*model.MatchRecallConfigure); ok {
			return NewMatcher(cfg)
		}
		zlog.LOG.Error("NewRecall.TypeAssertionFailed", zap.String("expected", "MatchRecallConfigure"), zap.String("actual", recallType))
	case model.RecallTypeModel:
		if cfg, ok := conf.(*model.ModelRecallConfigure); ok {
			return NewModeler(cfg)
		}
		zlog.LOG.Error("NewRecall.TypeAssertionFailed", zap.String("expected", "ModelRecallConfigure"), zap.String("actual", recallType))
	default:
		zlog.LOG.Error("NewRecall.UnknownType", zap.String("type", recallType))
	}

	return nil
}
