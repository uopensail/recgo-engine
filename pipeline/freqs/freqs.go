package freqs

import (
	"fmt"
	"sync"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IFilter defines an interface for filtering operations.
type IFilter interface {
	Do(userCtx *userctx.UserContext) model.IFilter
}

const (
	// userRecordKeyPrefix is the Redis (or feature storage) key prefix
	// for storing user action records.
	//
	// The full key format before splitting is:
	//   u_r_${action} -> "timestamp|itemKey"
	//
	// Example:
	//   u_r_click -> "1670832000|item123"
	//
	// This original string contains both a UNIX timestamp and an item ID/key
	// separated by actionSeparator ("|").
	//
	// During processing, this data is split into two separate fields:
	//   u_r_${action}_ts  -> stores only timestamps (int64 array)
	//   u_r_${action}_ids -> stores only item keys (string array)
	//
	// The splitting allows independent access to timestamps and item IDs,
	// which improves processing speed and simplifies filtering logic.
	userRecordKeyPrefix = "u_r_%s"

	// idsKeySuffix is the suffix appended to the base action key
	// to store only item IDs/keys after data splitting.
	// Example: "u_r_click_ids"
	idsKeySuffix = "_ids"

	// timestampKeySuffix is the suffix appended to the base action key
	// to store only timestamps after data splitting.
	// Example: "u_r_click_ts"
	timestampKeySuffix = "_ts"

	// actionSeparator defines the delimiter between timestamp and itemKey
	// in the original combined data format.
	// Example original data: "1670832000|item123"
	// After splitting by actionSeparator: ["1670832000", "item123"]
	actionSeparator = "|"
)

// FreqController controls frequency-based filtering logic.
type FreqController struct {
	frequencies   []model.IFreq // List of frequency rules
	enabledStatus []bool        // Enabled status for each frequency rule
}

// NewFreqController creates a new frequency controller and enables all frequencies by default.
func NewFreqController(frequencies []model.IFreq) *FreqController {
	pStat := prome.NewStat("NewFreqController")
	defer pStat.End()
	controller := &FreqController{
		frequencies:   frequencies,
		enabledStatus: make([]bool, len(frequencies)),
	}

	// Enable all frequencies initially
	for i := range controller.enabledStatus {
		controller.enabledStatus[i] = true
	}
	return controller
}

// Do executes all enabled frequency filters concurrently and merges their results.
func (fc *FreqController) Do(userCtx *userctx.UserContext) model.IFilter {
	pStat := prome.NewStat("FreqController.Do")
	defer pStat.End()

	var waitGroup sync.WaitGroup
	filterChan := make(chan *Filter, len(fc.frequencies))

	zlog.LOG.Debug("FreqController.Do", zap.Int("total_frequencies", len(fc.frequencies)))

	// Process each frequency rule in parallel
	for i := range fc.frequencies {
		if !fc.enabledStatus[i] {
			zlog.LOG.Info("FreqController.SkipFrequency", zap.Int("index", i))
			continue
		}

		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			filter := fc.processFrequency(userCtx, fc.frequencies[index])
			filterChan <- filter
		}(i)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()
	close(filterChan)

	// Merge all results into one filter
	resultFilter := NewFilter()
	for filter := range filterChan {
		if filter != nil {
			resultFilter.Merge(filter)
		}
	}

	zlog.LOG.Debug("FreqController.Do.Completed", zap.Int("final_filtered_count", len(resultFilter.ids)))
	return resultFilter
}

// processFrequency processes a single frequency rule and returns a Filter.
func (fc *FreqController) processFrequency(userCtx *userctx.UserContext, frequency model.IFreq) *Filter {
	actionKey := fmt.Sprintf(userRecordKeyPrefix, frequency.GetAction())

	// Fetch raw feature data
	rawKeys := userCtx.Features.Get(actionKey + idsKeySuffix)
	rawTimestamps := userCtx.Features.Get(actionKey + timestampKeySuffix)

	if rawKeys == nil || rawTimestamps == nil {
		zlog.LOG.Error("FreqController.ProcessFrequency.NoData", zap.String("actionKey", actionKey))
		return NewFilter()
	}

	itemKeys, err1 := rawKeys.GetStrings()
	timestamps, err2 := rawTimestamps.GetInt64s()

	if err1 != nil || err2 != nil || len(itemKeys) != len(timestamps) {
		zlog.LOG.Error("FreqController.ProcessFrequency.DataError",
			zap.Error(err1),
			zap.Error(err2),
			zap.Int("keys_len", len(itemKeys)),
			zap.Int("timestamps_len", len(timestamps)),
		)
		return NewFilter()
	}

	// Calculate frequencies
	frequencyMap := fc.calculateFrequency(itemKeys, timestamps, frequency.GetTimespan())

	// Create and return filter
	return fc.createFilter(userCtx, frequencyMap, frequency.GetFrequency())
}

// calculateFrequency counts occurrences of itemKeys within the given timespan.
func (fc *FreqController) calculateFrequency(itemKeys []string, timestamps []int64, timespan int) map[string]int {
	frequencyMap := make(map[string]int)
	cutoffTime := time.Now().Unix() - int64(timespan)

	for i, timestamp := range timestamps {
		if timestamp < cutoffTime {
			continue
		}
		itemKey := itemKeys[i]
		frequencyMap[itemKey]++
	}

	zlog.LOG.Debug("FreqController.CalculateFrequency.Result", zap.Int("unique_items", len(frequencyMap)))
	return frequencyMap
}

// createFilter creates a Filter from frequency data based on a threshold.
func (fc *FreqController) createFilter(userCtx *userctx.UserContext, frequencyMap map[string]int, threshold int) *Filter {
	filter := NewFilter()

	for itemKey, count := range frequencyMap {
		if count >= threshold {
			if itemID, _ := userCtx.Items.GetByKey(itemKey); itemID >= 0 {
				filter.Add(itemKey, itemID)
			} else {
				zlog.LOG.Error("FreqController.CreateFilter.InvalidItem", zap.String("itemKey", itemKey))
			}
		}
	}

	zlog.LOG.Debug("FreqController.CreateFilter.Completed", zap.Int("total_filtered", len(filter.ids)))
	return filter
}
