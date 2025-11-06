package report

import (
	"encoding/json"
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

type SLSLogReport struct {
	cfg config.SLSLogConfig
	p   *producer.Producer
}

func NewSLSLogReport(cfg config.SLSLogConfig) *SLSLogReport {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = cfg.Endpoint
	if len(cfg.RAM) > 0 {
		provider := sls.NewEcsRamRoleCredentialsProvider(cfg.RAM)
		producerConfig.CredentialsProvider = provider
	} else {
		provider := sls.NewStaticCredentialsProvider(cfg.AK,
			cfg.SK, "")
		producerConfig.CredentialsProvider = provider
	}

	producerInstance, err := producer.NewProducer(producerConfig)
	if err != nil {
		zlog.LOG.Error("SLSLogReport", zap.Error(err))
	}
	producerInstance.Start()
	return &SLSLogReport{p: producerInstance, cfg: cfg}
}

func (report *SLSLogReport) Report(uCtx *userctx.UserContext, resp *recapi.Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	contents := []*sls.LogContent{
		{
			Key:   proto.String("data"),       // 使用proto.String包装
			Value: proto.String(string(data)), // 转换为string并用proto.String包装
		},
	}

	err = report.p.SendLog(report.cfg.Project, report.cfg.LogStore, "rec_dist", os.Getenv("HOST"),
		&sls.Log{
			Time:     proto.Uint32(uint32(time.Now().Unix())),
			Contents: contents,
		})
	return err
}
func (report *SLSLogReport) Close() {
	if report.p != nil {
		report.p.SafeClose()
	}
}
