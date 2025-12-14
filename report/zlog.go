package report

import (
	"encoding/json"
	"fmt"

	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// ZLogReport implements IReport by logging recommendation results
// using the zlog logger. This is a simple, local-only reporter that does not
// send data to external systems.
//
// Useful for development or debugging environments.
type ZLogReport struct{}

// Report logs the recommendation response as a JSON string under the "data" key.
// Includes the user ID if available.
//
// Returns an error if JSON marshalling fails.
func (report *ZLogReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	if resp == nil {
		return fmt.Errorf("nil recommendation response")
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Log with optional user ID for better traceability
	fields := []zap.Field{zap.String("data", string(data))}
	if uCtx != nil && uCtx.Request.UserId != "" {
		fields = append(fields, zap.String("user_id", uCtx.Request.UserId))
	}

	zlog.LOG.Info("rec_dist", fields...)
	return nil
}

// Close is a no-op for ZLogReport since no external resources are used.
func (report *ZLogReport) Close() {}
