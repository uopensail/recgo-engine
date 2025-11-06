package strategy

import (
	"fmt"

	"github.com/spaolacci/murmur3"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/pipeline"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/sample"
)

type Strategy struct {
	feeds           map[string]pipeline.IPipeline
	related         map[string]pipeline.IPipeline
	feedsBuckets    []pipeline.IPipeline
	releatedBuckets []pipeline.IPipeline
}

func NewStrategy(conf *config.AppConfig) *Strategy {
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

	feedsBuckets := make([]pipeline.IPipeline, 100)
	for _, p := range feeds {
		for _, idx := range p.GetBuckets() {
			feedsBuckets[idx] = p
		}
	}
	releatedBuckets := make([]pipeline.IPipeline, 100)
	for _, p := range related {
		for _, idx := range p.GetBuckets() {
			releatedBuckets[idx] = p
		}
	}

	if len(feeds) == 0 || len(related) == 0 {
		panic(fmt.Errorf("build strategy fail"))
	}

	return &Strategy{feeds: feeds, related: related, feedsBuckets: feedsBuckets, releatedBuckets: releatedBuckets}
}

func (s *Strategy) Feeds(uCtx *userctx.UserContext) *recapi.Response {
	h := murmur3.New64()
	h.Write([]byte(uCtx.Request.UserId))
	p := s.feedsBuckets[h.Sum64()%100]
	if p == nil {
		panic(fmt.Errorf("no pipeline hit"))
	}
	collection := p.Do(uCtx)
	resp := &recapi.Response{
		TraceId:  uCtx.Request.TraceId,
		UserId:   uCtx.Request.UserId,
		Pipeline: p.GetName(),
	}
	resp.Items = make([]*recapi.ItemInfo, 0, len(collection))
	var fea sample.Feature
	for _, entry := range collection {
		fea, _ = entry.Get(model.ChannelsKey)
		channels, _ := fea.GetStrings()
		fea, _ = entry.Get(model.ReasonsKey)
		reasopns, _ := fea.GetStrings()

		resp.Items = append(resp.Items, &recapi.ItemInfo{
			Item:     entry.Key,
			Channels: channels,
			Reasons:  reasopns,
		})
	}
	return resp
}

func (s *Strategy) Related(uCtx *userctx.UserContext) *recapi.Response {
	h := murmur3.New64()
	h.Write([]byte(uCtx.Request.UserId))
	p := s.releatedBuckets[h.Sum64()%100]
	if p == nil {
		panic(fmt.Errorf("no pipeline hit"))
	}
	collection := p.Do(uCtx)
	resp := &recapi.Response{
		TraceId:  uCtx.Request.TraceId,
		UserId:   uCtx.Request.UserId,
		Pipeline: p.GetName(),
	}
	resp.Items = make([]*recapi.ItemInfo, 0, len(collection))
	var fea sample.Feature
	for _, entry := range collection {
		fea, _ = entry.Get(model.ChannelsKey)
		channels, _ := fea.GetStrings()
		fea, _ = entry.Get(model.ReasonsKey)
		reasopns, _ := fea.GetStrings()

		resp.Items = append(resp.Items, &recapi.ItemInfo{
			Item:     entry.Key,
			Channels: channels,
			Reasons:  reasopns,
		})
	}
	return resp
}

var StrategyInstance *Strategy
