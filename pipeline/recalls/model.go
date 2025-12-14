package recalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Request represents the payload sent to a model recall service.
// Features: user context features.
// Excludes: keys to be excluded from recall results.
// Count:    number of items requested.
type Request struct {
	Features sample.Features `json:"features"`
	Excludes []string        `json:"excludes"`
	Count    int             `json:"count"`
}

// Response represents the payload returned by a model recall service.
// Code:    status code from the model service.
// Msg:     optional status/error message.
// Data:    list of recalled key-score pairs.
type Response struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []model.KeyScore `json:"data"`
}

// Modeler is an implementation of IRecall that fetches recall results from an external model service via HTTP.
type Modeler struct {
	conf   *model.ModelRecallConfigure // recall configuration, contains URL & other parameters
	client *http.Client                // HTTP client
}

// NewModeler creates a Modeler with customized HTTP client settings.
func NewModeler(conf *model.ModelRecallConfigure) *Modeler {
	pStat := prome.NewStat("Recall.NewModeler")
	defer pStat.End()

	client := &http.Client{
		Timeout: 30 * time.Second, // request timeout
		Transport: &http.Transport{
			MaxIdleConns:        100,              // max idle connections globally
			MaxIdleConnsPerHost: 10,               // max idle connections per host
			IdleConnTimeout:     90 * time.Second, // idle connection timeout
		},
	}

	zlog.LOG.Info("Modeler.Created",
		zap.String("name", conf.Name),
		zap.String("url", conf.URL),
		zap.Int("count", conf.Count),
	)

	return &Modeler{
		conf:   conf,
		client: client,
	}
}

// Do executes the model-based recall process:
//  1. Prepare request payload.
//  2. Send POST to model recall endpoint.
//  3. Parse JSON response.
//  4. Convert response to model.Collection.
func (m *Modeler) Do(uCtx *userctx.UserContext) model.Collection {
	pStat := prome.NewStat("Recall.Modeler.Do")
	defer pStat.End()

	// Prepare request payload
	request := Request{
		Features: uCtx.Features,
		Excludes: uCtx.Filter.Exclude(),
		Count:    m.conf.Count,
	}
	data, err := json.Marshal(request)
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Modeler.Do.MarshalError", zap.Error(err))
		return nil
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", m.conf.URL, bytes.NewBuffer(data))
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Modeler.Do.NewRequestError", zap.Error(err))
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send HTTP request
	response, err := m.client.Do(req)
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Modeler.Do.HTTPError", zap.Error(err))
		return nil
	}
	defer response.Body.Close()

	// Check for non-200 status codes
	if response.StatusCode != http.StatusOK {
		pStat.MarkErr()
		zlog.LOG.Warn("Recall.Modeler.Do.NonOKStatus",
			zap.Int("status_code", response.StatusCode),
			zap.String("status", response.Status))
		return nil
	}

	// Read body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Modeler.Do.ReadBodyError", zap.Error(err))
		return nil
	}

	// Parse JSON response
	resp := Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Modeler.Do.UnmarshalError", zap.Error(err))
		return nil
	}

	// Check for application-level errors in code/msg
	if resp.Code != 0 {
		pStat.MarkErr()
		zlog.LOG.Warn("Recall.Modeler.Do.AppError",
			zap.Int("code", resp.Code),
			zap.String("msg", resp.Msg))
		return nil
	}

	// Build collection from response data
	collection := make([]*model.Entry, 0, len(resp.Data))
	for _, kv := range resp.Data {
		id, feas := uCtx.Items.GetByKey(kv.Key)
		if id < 0 {
			continue
		}

		entry := &model.Entry{
			ID:       id,
			KeyScore: kv,
			Runtime:  *model.NewRuntime(feas),
		}
		entry.AddChan(m.conf.Name, fmt.Sprintf("recall by model: %s", m.conf.URL))
		collection = append(collection, entry)
	}

	count := minInt(m.conf.Count, len(collection))
	zlog.LOG.Debug("Recall.Modeler.Do.Completed",
		zap.Int("requested_count", m.conf.Count),
		zap.Int("returned_count", count),
	)
	pStat.SetCounter(count)

	return collection[:count]
}
