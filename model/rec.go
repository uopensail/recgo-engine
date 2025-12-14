package model

import (
	"fmt"

	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// IFilter defines a filter interface used in recommendation pipeline.
// Exists returns true if the ID should be filtered out.
// Exclude returns a list of keys to be excluded.
type IFilter interface {
	Exists(id int) bool // True: filtered, False: pass
	Exclude() []string
}

// Resource defines basic information for loaded resources.
type Resource interface {
	GetUpdateTime() int64
	GetURL() string
}

// Keys used in runtime features to store channels and reasons.
const (
	ChannelsKey = "i_ctx_chans"
	ReasonsKey  = "i_ctx_reasons"
)

// Entry represents a candidate item in recommendation pipeline,
// including its item ID, score, and runtime features (channels, reasons, etc.).
type Entry struct {
	ID int
	KeyScore
	Runtime
}

// NewEntry creates a new Entry from a KeyScore and Items.
// It initializes empty channel and reason lists in runtime features.
// Returns an error if the key cannot be found in Items.
func NewEntry(k KeyScore, items *Items) (*Entry, error) {
	id, feas := items.GetByKey(k.Key)
	if feas == nil {
		zlog.LOG.Error("Entry.NewEntry.KeyNotFound", zap.String("key", k.Key))
		return nil, fmt.Errorf("key miss: %s", k.Key)
	}

	r := NewRuntime(feas)
	r.Set(ChannelsKey, &sample.Strings{Value: make([]string, 0, 8)})
	r.Set(ReasonsKey, &sample.Strings{Value: make([]string, 0, 8)})

	zlog.LOG.Info("Entry.NewEntry.Created", zap.Int("id", id), zap.String("key", k.Key))
	return &Entry{id, k, *r}, nil
}

// AddChan adds a channel and reason to the entry's runtime features.
// This will update both channel and reason lists.
func (entry *Entry) AddChan(channel string, reason string) {
	// Update channels
	feaChan, _ := entry.Get(ChannelsKey)
	chList, _ := feaChan.GetStrings()
	chList = append(chList, channel)
	entry.Set(ChannelsKey, &sample.Strings{Value: chList})

	// Update reasons
	feaReason, _ := entry.Get(ReasonsKey)
	reasonList, _ := feaReason.GetStrings()
	reasonList = append(reasonList, reason)
	entry.Set(ReasonsKey, &sample.Strings{Value: reasonList})

	zlog.LOG.Info("Entry.AddChan", zap.Int("id", entry.ID), zap.String("channel", channel), zap.String("reason", reason))
}

// MergeChans merges channels and reasons from another Entry into this Entry.
// The source Entry is not modified.
func (entry *Entry) MergeChans(src *Entry) {
	// Merge channels
	feaChanDst, _ := entry.Get(ChannelsKey)
	dstChList, _ := feaChanDst.GetStrings()

	feaChanSrc, _ := src.Get(ChannelsKey)
	srcChList, _ := feaChanSrc.GetStrings()

	dstChList = append(dstChList, srcChList...)
	entry.Set(ChannelsKey, &sample.Strings{Value: dstChList})

	// Merge reasons
	feaReasonDst, _ := entry.Get(ReasonsKey)
	dstReasonList, _ := feaReasonDst.GetStrings()

	feaReasonSrc, _ := src.Get(ReasonsKey)
	srcReasonList, _ := feaReasonSrc.GetStrings()

	dstReasonList = append(dstReasonList, srcReasonList...)
	entry.Set(ReasonsKey, &sample.Strings{Value: dstReasonList})

	zlog.LOG.Info("Entry.MergeChans",
		zap.Int("dst_id", entry.ID),
		zap.Int("src_id", src.ID),
		zap.Int("merged_channels_count", len(srcChList)),
		zap.Int("merged_reasons_count", len(srcReasonList)),
	)
}

// Collection is a slice of Entry pointers.
// It implements sort.Interface to allow sorting by Score in descending order.
type Collection []*Entry

// Less returns true if entry at index i has a higher score than entry at index j.
func (c Collection) Less(i, j int) bool {
	return c[i].KeyScore.Score > c[j].KeyScore.Score
}

// Len returns the number of entries in the collection.
func (c Collection) Len() int {
	return len(c)
}

// Swap exchanges the entries at indexes i and j.
func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
