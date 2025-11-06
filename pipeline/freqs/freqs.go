package freqs

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
)

type IFilter interface {
	Do(userCtx *userctx.UserContext) model.IFliter
}

const (
	userRecordKeyPrefix = "u_r_%s"
	idsKeySuffix        = "_ids"
	timestampKeySuffix  = "_ts"
	actionSeparator     = "|"
)

type FreqController struct {
	frequencies   []model.IFreq
	enabledStatus []bool
}

func NewFreqController(frequencies []model.IFreq) *FreqController {
	controller := &FreqController{
		frequencies:   frequencies,
		enabledStatus: make([]bool, len(frequencies)),
	}

	// 初始化所有频率控制器为启用状态
	for i := range controller.enabledStatus {
		controller.enabledStatus[i] = true
	}

	return controller
}

func (fc *FreqController) Do(userCtx *userctx.UserContext) model.IFliter {
	var waitGroup sync.WaitGroup
	filterChan := make(chan *Filter, len(fc.frequencies))

	// 并发处理每个频率控制器
	for i := range fc.frequencies {
		if !fc.enabledStatus[i] {
			continue
		}

		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			filter := fc.processFrequency(userCtx, fc.frequencies[index])
			filterChan <- filter
		}(i)
	}

	// 等待所有goroutine完成后再关闭channel

	waitGroup.Wait()
	close(filterChan)

	// 合并所有过滤结果
	resultFilter := NewFilter()
	for filter := range filterChan {
		if filter != nil {
			resultFilter.Merge(filter)
		}
	}

	return resultFilter
}

func (fc *FreqController) parseActionData(actions []string) ([]string, []int64) {
	itemKeys := make([]string, 0, len(actions))
	timestamps := make([]int64, 0, len(actions))

	for _, action := range actions {
		parts := strings.Split(action, actionSeparator)
		if len(parts) != 2 {
			continue // 跳过格式不正确的数据
		}

		timestamp, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue // 跳过时间戳解析失败的数据
		}

		itemKeys = append(itemKeys, parts[1])
		timestamps = append(timestamps, timestamp)
	}

	return itemKeys, timestamps
}

func (fc *FreqController) processFrequency(userCtx *userctx.UserContext, frequency model.IFreq) *Filter {
	actionKey := fmt.Sprintf(userRecordKeyPrefix, frequency.GetAction())

	// 获取解析后的数据
	rawKeys := userCtx.Features.Get(actionKey + idsKeySuffix)
	rawTimestamps := userCtx.Features.Get(actionKey + timestampKeySuffix)

	if rawKeys == nil || rawTimestamps == nil {
		return NewFilter()
	}

	itemKeys, err1 := rawKeys.GetStrings()
	timestamps, err2 := rawTimestamps.GetInt64s()

	if err1 != nil || err2 != nil || len(itemKeys) != len(timestamps) {
		return NewFilter()
	}

	// 计算频率统计
	frequencyMap := fc.calculateFrequency(itemKeys, timestamps, frequency.GetTimespan())

	// 生成过滤器
	return fc.createFilter(userCtx, frequencyMap, frequency.GetFrequency())
}

func (fc *FreqController) calculateFrequency(itemKeys []string, timestamps []int64, timespan int) map[string]int {
	frequencyMap := make(map[string]int)
	cutoffTime := time.Now().Unix() - int64(timespan)

	for i, timestamp := range timestamps {
		if timestamp < cutoffTime {
			continue // 超出时间窗口的数据不计算
		}

		itemKey := itemKeys[i]
		frequencyMap[itemKey]++
	}

	return frequencyMap
}

func (fc *FreqController) createFilter(userCtx *userctx.UserContext, frequencyMap map[string]int, threshold int) *Filter {
	filter := NewFilter()

	for itemKey, count := range frequencyMap {
		if count >= threshold {
			if itemID, err := userCtx.Items.GetByKey(itemKey); err == nil && itemID >= 0 {
				filter.Add(itemKey, itemID)
			}
		}
	}

	return filter
}
