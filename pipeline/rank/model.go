package rank

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Request represents the payload sent to the model-based ranking service.
type Request struct {
	Features   sample.Features `json:"features"`
	Collection []string        `json:"collection"`
}

// Response represents the ranking results returned by the model service.
type Response struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []model.KeyScore `json:"data"`
}

// Modeler ranks items using scores from an external model service.
type Modeler struct {
	conf   *model.ModelBasedRankConfigure // rank configuration
	client *http.Client                   // HTTP client
}

// NewModeler creates a new Modeler with configured HTTP client settings.
func NewModeler(conf *model.ModelBasedRankConfigure) *Modeler {
	pStat := prome.NewStat("Rank.NewModeler")
	defer pStat.End()
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	return &Modeler{
		conf:   conf,
		client: client,
	}
}

// Do sends the current collection to the model service for scoring,
// updates each entry's Score, and returns the sorted collection.
func (m *Modeler) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("Rank.Modeler.Do")
	defer pStat.End()
	// Prepare list of keys
	keys := make([]string, 0, len(collection))
	for _, entry := range collection {
		keys = append(keys, entry.Key)
	}

	// Construct request payload
	request := Request{
		Features:   uCtx.Features,
		Collection: keys,
	}
	data, err := json.Marshal(request)
	if err != nil {
		zlog.LOG.Error("Rank.Modeler.Do.MarshalError", zap.Error(err))
		return collection
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", m.conf.URL, bytes.NewBuffer(data))
	if err != nil {
		zlog.LOG.Error("Rank.Modeler.Do.NewRequestError", zap.Error(err))
		return collection
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send HTTP request
	resp, err := m.client.Do(req)
	if err != nil {
		zlog.LOG.Error("Rank.Modeler.Do.HTTPError", zap.Error(err))
		return collection
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		zlog.LOG.Warn("Rank.Modeler.Do.NonOKStatus",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status))
		return collection
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zlog.LOG.Error("Rank.Modeler.Do.ReadBodyError", zap.Error(err))
		return collection
	}
	var responseData Response
	if err := json.Unmarshal(body, &responseData); err != nil {
		zlog.LOG.Error("Rank.Modeler.Do.UnmarshalError", zap.Error(err))
		return collection
	}

	// Business-level error check
	if responseData.Code != 0 {
		zlog.LOG.Warn("Rank.Modeler.Do.AppError",
			zap.Int("code", responseData.Code),
			zap.String("msg", responseData.Msg))
		return collection
	}

	// Map scores
	scoreMap := make(map[string]float32, len(responseData.Data))
	for _, kv := range responseData.Data {
		scoreMap[kv.Key] = kv.Score
	}

	// Assign scores to collection
	for _, entry := range collection {
		if score, ok := scoreMap[entry.Key]; ok {
			entry.KeyScore.Score = score
		} else {
			entry.KeyScore.Score = 0
		}
	}

	// Sort using the built-in sort.Interface of model.Collection
	sort.Stable(collection)

	zlog.LOG.Debug("Rank.Modeler.Do.Completed",
		zap.Int("total", len(collection)),
		zap.String("url", m.conf.URL))
	return collection
}
