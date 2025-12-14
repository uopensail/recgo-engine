package model

import (
	"bufio"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/bytedance/sonic"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Items represents a collection of immutable feature data loaded from file.
// It supports lookup by key or ID, and stores features in memory for fast access.
type Items struct {
	arena      *sample.Arena               // Memory arena for ImmutableFeatures
	dict       map[string]int              // Map from item key to index in array
	array      []*sample.ImmutableFeatures // Immutable feature list
	filePath   string                      // Source file path
	updateTime int64                       // UNIX timestamp when data was last updated
}

// NewItems loads Items from a tab-delimited file.
// File format: each line contains "<key>\t<json>"
// Example:
//
//	item123   {"feature1":{...},"feature2":{...}}
//
// Each entry's JSON is parsed into ImmutableFeatures stored in memory.
// Logs:
// - Error if file cannot be opened
// - Warning if a line is skipped due to format error
// - Error if JSON parsing fails
// - Info on total items loaded and time taken
func NewItems(filePath string) (Resource, error) {
	stat := prome.NewStat("NewItems")
	defer stat.End()

	startTime := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		zlog.LOG.Error("Items.FileOpenError", zap.String("filePath", filePath), zap.Error(err))
		stat.MarkErr()
		return nil, err
	}
	defer file.Close()
	zlog.LOG.Info("Items.FileOpenSuccess", zap.String("filePath", filePath))

	scanner := bufio.NewScanner(file)
	items := &Items{
		arena:    sample.NewArena(),
		array:    make([]*sample.ImmutableFeatures, 0, 4096),
		dict:     make(map[string]int, 4096),
		filePath: filePath,
	}

	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Split(line, "\t")
		if len(ss) != 2 {
			zlog.LOG.Warn("Items.SkipLine.InvalidFormat", zap.Int("line_index", index), zap.String("line", line))
			continue
		}

		feas := sample.NewImmutableFeatures(items.arena)
		err = sonic.Unmarshal(unsafe.Slice(unsafe.StringData(ss[1]), len(ss[1])), feas)
		if err != nil {
			zlog.LOG.Error("Items.JSONUnmarshalError", zap.String("key", ss[0]), zap.String("raw_data", ss[1]), zap.Error(err))
			continue
		}

		items.array = append(items.array, feas)
		items.dict[ss[0]] = index
		index++
	}

	if err := scanner.Err(); err != nil {
		zlog.LOG.Error("Items.ScannerError", zap.Error(err))
		stat.MarkErr()
		return nil, err
	}

	// Log stats and update time
	items.updateTime = time.Now().Unix()
	stat.SetCounter(index)

	zlog.LOG.Info("Items.LoadComplete",
		zap.Int("total_items", index),
		zap.Duration("elapsed", time.Since(startTime)),
	)

	return items, nil
}

// GetByKey retrieves the ID and ImmutableFeatures for a given key.
// Returns -1 and nil if the key does not exist.
func (items *Items) GetByKey(key string) (int, *sample.ImmutableFeatures) {
	if id, ok := items.dict[key]; ok {
		return id, items.array[id]
	}
	return -1, nil
}

// GetByID retrieves the ImmutableFeatures for a given ID.
// Returns nil if the ID is out of range.
func (items *Items) GetByID(id int) *sample.ImmutableFeatures {
	if id >= 0 && id < len(items.array) {
		return items.array[id]
	}
	return nil
}

// GetUpdateTime returns the UNIX timestamp of the last data update.
func (items *Items) GetUpdateTime() int64 {
	return items.updateTime
}

// GetURL returns the source file path of the Items data.
func (items *Items) GetURL() string {
	return items.filePath
}
