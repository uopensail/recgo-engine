package strategy

import (
	"fmt"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/pipeline"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
)

// Strategy routes user requests to specific pipelines based on pipeline name in the request.
type Strategy struct {
	feeds   map[string]pipeline.IPipeline
	related map[string]pipeline.IPipeline
}

// NewStrategy builds a new Strategy from AppConfig, initializing pipelines for feeds and related.
func NewStrategy(conf *config.AppConfig) *Strategy {
	pStat := prome.NewStat("NewStrategy")
	defer pStat.End()

	feeds := make(map[string]pipeline.IPipeline, len(conf.Feeds))
	for _, pconf := range conf.Feeds {
		p := pipeline.NewPipeline(&pconf)
		if p != nil {
			feeds[pconf.Name] = p
		}
	}

	related := make(map[string]pipeline.IPipeline, len(conf.Related))
	for _, pconf := range conf.Related {
		p := pipeline.NewPipeline(&pconf)
		if p != nil {
			related[pconf.Name] = p
		}
	}

	if len(feeds) == 0 || len(related) == 0 {
		panic(fmt.Errorf("build strategy fail: feeds or related pipelines missing"))
	}

	return &Strategy{
		feeds:   feeds,
		related: related,
	}
}

// runPipeline executes the given pipeline for the given user context and builds a standard Response.
func (s *Strategy) runPipeline(uCtx *userctx.UserContext, p pipeline.IPipeline) *recapi.Response {
	if p == nil {
		// Pipeline not found, return error response
		return &recapi.Response{
			Code:     -1,
			Message:  fmt.Sprintf("pipeline not found: %s", uCtx.Request.Pipeline),
			TraceId:  uCtx.Request.TraceId,
			UserId:   uCtx.Request.UserId,
			Pipeline: "",
			Items:    []*recapi.ItemInfo{},
			Count:    0,
		}
	}

	collection := p.Do(uCtx)
	resp := &recapi.Response{
		Code:     0,
		Message:  "success",
		TraceId:  uCtx.Request.TraceId,
		UserId:   uCtx.Request.UserId,
		Pipeline: p.GetName(),
		Items:    make([]*recapi.ItemInfo, 0, len(collection)),
		Count:    len(collection),
	}

	var fea sample.Feature
	for _, entry := range collection {
		fea, _ = entry.Get(model.ChannelsKey)
		channels, _ := fea.GetStrings()
		fea, _ = entry.Get(model.ReasonsKey)
		reasons, _ := fea.GetStrings()

		resp.Items = append(resp.Items, &recapi.ItemInfo{
			Item:     entry.Key,
			Channels: channels,
			Reasons:  reasons,
		})
	}

	return resp
}

// Feeds returns feed recommendations for the given user context.
func (s *Strategy) Feeds(uCtx *userctx.UserContext) *recapi.Response {
	pStat := prome.NewStat(fmt.Sprintf("Strategy.Feed.%s", uCtx.Request.Pipeline))
	defer pStat.End()

	if p, ok := s.feeds[uCtx.Request.Pipeline]; ok {
		return s.runPipeline(uCtx, p)
	}
	pStat.MarkErr()
	return s.runPipeline(uCtx, nil)
}

// Related returns related item recommendations for the given user context.
func (s *Strategy) Related(uCtx *userctx.UserContext) *recapi.Response {
	pStat := prome.NewStat(fmt.Sprintf("Strategy.Related.%s", uCtx.Request.Pipeline))
	defer pStat.End()

	if p, ok := s.related[uCtx.Request.Pipeline]; ok {
		return s.runPipeline(uCtx, p)
	}
	pStat.MarkErr()
	return s.runPipeline(uCtx, nil)
}

// StrategyInstance is the global strategy singleton.
var StrategyInstance *Strategy
