package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// SLSLogReport implements IReport by sending recommendation logs
// to Alibaba Cloud's Simple Log Service (SLS) using an asynchronous producer.
type SLSLogReport struct {
	cfg config.SLSLogConfig
	p   *producer.Producer
}

// NewSLSLogReport creates a new SLSLogReport using the provided configuration.
// Depending on config.RAM, it uses either RAM role credentials or static AK/SK.
// The producer is started automatically.
//
// Returns nil if the producer cannot be created.
//
// Example:
//
//	cfg := config.SLSLogConfig{Endpoint: "...", AK: "...", SK: "...", Project: "...", LogStore: "..."}
//	report := report.NewSLSLogReport(cfg)
//	defer report.Close()
func NewSLSLogReport(cfg config.SLSLogConfig) *SLSLogReport {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = cfg.Endpoint

	// Set credentials provider
	if cfg.RAM != "" {
		producerConfig.CredentialsProvider = sls.NewEcsRamRoleCredentialsProvider(cfg.RAM)
	} else {
		producerConfig.CredentialsProvider = sls.NewStaticCredentialsProvider(cfg.AK, cfg.SK, "")
	}

	// Create producer
	producerInstance, err := producer.NewProducer(producerConfig)
	if err != nil {
		zlog.LOG.Error("failed to create SLS producer", zap.Error(err))
		return nil
	}

	producerInstance.Start()
	return &SLSLogReport{p: producerInstance, cfg: cfg}
}

// Report sends the recommendation response to Alibaba Cloud SLS.
// The data is marshaled into JSON and sent with the key "data".
// The LogStore and Project come from the config.
// Host is retrieved from the HOST environment variable.
//
// Returns an error if marshalling fails or sending the log fails.
func (report *SLSLogReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	if report == nil || report.p == nil {
		return fmt.Errorf("SLSLogReport is not initialized")
	}

	if resp == nil {
		return fmt.Errorf("nil recommendation response")
	}

	// Marshal response
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Build log contents
	contents := []*sls.LogContent{
		{
			Key:   proto.String("data"),
			Value: proto.String(string(data)),
		},
	}

	// Send log
	return report.p.SendLog(report.cfg.Project, report.cfg.LogStore, "rec_dist", os.Getenv("HOST"),
		&sls.Log{
			Time:     proto.Uint32(uint32(time.Now().Unix())),
			Contents: contents,
		})
}

// Close gracefully shuts down the SLS producer, ensuring any pending logs are flushed.
func (report *SLSLogReport) Close() {
	if report.p != nil {
		report.p.SafeClose()
	}
}
