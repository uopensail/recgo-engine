package rank

import (
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
)

type ChannelPriority struct {
	conf *model.ChannelPriorityRankConfigure
}

func NewChannelPriority(conf *model.ChannelPriorityRankConfigure) *ChannelPriority {
	conf.Priorities = append(conf.Priorities, "")
	return &ChannelPriority{
		conf: conf,
	}
}

func (cp *ChannelPriority) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	channels := make(map[string]model.Collection, len(cp.conf.Priorities)+1)
	for _, ch := range cp.conf.Priorities {
		channels[ch] = make(model.Collection, 0, 64)
	}

	remain, _ := channels[""]

	for _, entry := range collection {
		chans, _ := entry.Get(model.ChannelsKey)
		chs, _ := chans.GetStrings()
		for _, ch := range chs {
			if arr, ok := channels[ch]; ok {
				arr = append(arr, entry)
			} else {
				remain = append(remain, entry)
			}
		}
	}

	filter := make(map[int]struct{})
	ret := make(model.Collection, 0, len(collection))
	for _, ch := range cp.conf.Priorities {
		for _, entry := range channels[ch] {
			if _, ok := filter[entry.ID]; ok {
				continue
			} else {
				filter[entry.ID] = struct{}{}
				ret = append(ret, entry)
			}
		}
	}
	return ret
}
