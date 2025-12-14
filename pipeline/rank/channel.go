package rank

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// ChannelPriority reorders collection entries based on configured channel priorities.
// Entries belonging to higher priority channels appear first in the output.
type ChannelPriority struct {
	conf *model.ChannelPriorityRankConfigure // sorting configuration
}

// NewChannelPriority creates a new ChannelPriority ranker.
// It ensures there is an empty string priority representing the "remain" / unmatched channel group.
func NewChannelPriority(conf *model.ChannelPriorityRankConfigure) *ChannelPriority {
	pStat := prome.NewStat("NewChannelPriority")
	defer pStat.End()
	conf.Priorities = append(conf.Priorities, "") // add remain group
	return &ChannelPriority{
		conf: conf,
	}
}

// Do reorders the given collection based on channel priorities.
// 1. Group entries into channels according to entry's ChannelsKey feature.
// 2. Traverse priorities, appending entries from each channel to result (without duplicates).
// 3. Append remaining unmatched entries at the end.
// Returns a new ordered collection.
func (cp *ChannelPriority) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("ChannelPriority.Do")
	defer pStat.End()
	// Initialize channel groups
	channels := make(map[string]model.Collection, len(cp.conf.Priorities))
	for _, ch := range cp.conf.Priorities {
		channels[ch] = make(model.Collection, 0, 64)
	}

	// Group entries by their channels
	for _, entry := range collection {
		chans, _ := entry.Get(model.ChannelsKey)
		chs, _ := chans.GetStrings()
		inserted := false
		for _, ch := range chs {
			if arr, ok := channels[ch]; ok {
				channels[ch] = append(arr, entry)
				inserted = true
			}
		}
		if !inserted {
			// Put into remain (empty string channel)
			channels[""] = append(channels[""], entry)
		}
	}

	// Merge channels by priority order, ensuring no duplicates
	filter := make(map[int]struct{}) // track appended entry IDs
	ret := make(model.Collection, 0, len(collection))
	for _, ch := range cp.conf.Priorities {
		for _, entry := range channels[ch] {
			if _, exists := filter[entry.ID]; exists {
				continue
			}
			filter[entry.ID] = struct{}{}
			ret = append(ret, entry)
		}
	}

	zlog.LOG.Debug("ChannelPriority.Do.Completed",
		zap.Int("total_returned", len(ret)),
		zap.Any("priorities", cp.conf.Priorities),
	)
	return ret
}
