package userctx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Request defines the feature fetch request payload for user context creation.
type Request struct {
	UserId string `json:"user_id"`
}

// Response defines the structure of the feature fetch API response.
type Response struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *sample.MutableFeatures `json:"data"`
}

// UserContext holds all runtime information for recommendation processing.
// It contains request info, loaded items, filter results, user features and related features (contextual items).
type UserContext struct {
	context.Context
	Request  *recapi.Request
	Items    *model.Items
	Filter   model.IFilter
	Features *sample.MutableFeatures
	Related  *sample.ImmutableFeatures
}

// NewUserContext creates a UserContext from a base context and recommendation API request.
// It loads item resources, fetches related item features if RelateId is provided,
// merges request-level features into the context, and fetches remote user features if available.
func NewUserContext(ctx context.Context, req *recapi.Request) *UserContext {
	pStat := prome.NewStat("NewUserContext")
	defer pStat.End()

	items := resources.ResourceManagerInstance.GetItems()

	// Load related item features if RelateId is provided
	var related *sample.ImmutableFeatures
	if len(req.RelateId) > 0 {
		id, feas := items.GetByKey(req.RelateId)
		if id >= 0 {
			related = feas
		}
	}

	uCtx := UserContext{
		Context:  ctx,
		Request:  req,
		Items:    items,
		Filter:   nil,
		Features: nil,
		Related:  related,
	}

	// Fetch remote user features
	features := uCtx.fetchFeatures()
	if features == nil {
		zlog.LOG.Warn("UserContext.FetchFeaturesEmpty",
			zap.String("user_id", req.UserId))
		features = sample.NewMutableFeatures()
	}

	// Merge request features and context features into the user feature set
	merge := func(key string, feature sample.Feature) error {
		features.Set(key, feature)
		return nil
	}
	uCtx.Request.Features.ForEach(merge)
	uCtx.Request.Context.ForEach(merge)
	uCtx.Features = features

	zlog.LOG.Debug("UserContext.Created",
		zap.String("user_id", req.UserId),
		zap.Int("feature_count", features.Len()))

	return &uCtx
}

// fetchFeatures contacts the external user feature service to retrieve user features.
// Returns nil if the request fails or the response is invalid.
func (uCtx *UserContext) fetchFeatures() *sample.MutableFeatures {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	reqPayload := Request{
		UserId: uCtx.Request.UserId,
	}
	data, err := json.Marshal(reqPayload)
	if err != nil {
		zlog.LOG.Error("UserContext.MarshalRequestFailed",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Error(err))
		return nil
	}

	// TODO: Replace "/user" with a full URL or configurable endpoint
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	if err != nil {
		zlog.LOG.Error("UserContext.CreateHTTPRequestFailed",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Error(err))
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := client.Do(req)
	if err != nil {
		zlog.LOG.Error("UserContext.FetchFeatureHTTPRequestFailed",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Error(err))
		return nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.LOG.Error("UserContext.ReadResponseBodyFailed",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Error(err))
		return nil
	}

	resp := Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		zlog.LOG.Error("UserContext.UnmarshalResponseFailed",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Error(err))
		return nil
	}
	if resp.Data == nil {
		zlog.LOG.Warn("UserContext.EmptyFeatureResponse",
			zap.String("user_id", uCtx.Request.UserId),
			zap.Int("code", resp.Code),
			zap.String("msg", resp.Msg))
	}

	return resp.Data
}
