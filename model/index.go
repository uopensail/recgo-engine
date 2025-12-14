package model

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// KeyScore represents an item with its relevance score.
type KeyScore struct {
	Key   string  `json:"key"`   // Unique identifier of the candidate item
	Score float32 `json:"score"` // Relevance score for the candidate item
}

// IndexEntry represents a single entry in the inverted index.
type IndexEntry struct {
	Key    string     `json:"key"`    // Index key (feature value)
	Values []KeyScore `json:"values"` // List of candidates associated with this key
}

// InvertedIndex stores an inverted index for fast candidate retrieval.
type InvertedIndex struct {
	indexMap   map[string]IndexEntry // Internal map for O(1) key lookup
	filePath   string                // Filepath of inverted index source
	updateTime int64                 // UNIX timestamp when the index was last updated
}

// NewInvertedIndex loads an InvertedIndex instance from a JSON file.
// The JSON file should contain an array of IndexEntry objects.
//
// Example JSON:
//
//	[
//	  {
//	    "key": "feature1",
//	    "values": [
//	      {"key": "itemA", "score": 0.95},
//	      {"key": "itemB", "score": 0.87}
//	    ]
//	  }
//	]
//
// Candidates within each entry will be sorted by score in descending order.
//
// Logs:
// - Info on file read success/failure
// - Info on JSON parsing result
// - Info on total index entries loaded
func NewInvertedIndex(filePath string) (Resource, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		zlog.LOG.Info("InvertedIndex.FileReadError", zap.String("filePath", filePath), zap.Error(err))
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	zlog.LOG.Info("InvertedIndex.FileReadSuccess", zap.String("filePath", filePath), zap.Int("bytes_read", len(data)))

	// Parse JSON data
	var entries []IndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		zlog.LOG.Info("InvertedIndex.JSONParseError", zap.String("filePath", filePath), zap.Error(err))
		return nil, fmt.Errorf("failed to parse JSON from file %s: %w", filePath, err)
	}
	zlog.LOG.Info("InvertedIndex.JSONParseSuccess", zap.Int("entries_count", len(entries)))

	// Build index map
	indexMap := make(map[string]IndexEntry, len(entries))
	for _, entry := range entries {
		// Sort candidates by score in descending order
		sort.Slice(entry.Values, func(i, j int) bool {
			return entry.Values[i].Score > entry.Values[j].Score
		})
		indexMap[entry.Key] = entry
	}

	zlog.LOG.Info("InvertedIndex.BuildComplete", zap.Int("total_keys", len(indexMap)))

	return &InvertedIndex{
		indexMap:   indexMap,
		filePath:   filePath,
		updateTime: time.Now().Unix(),
	}, nil
}

// Get retrieves the IndexEntry for the given key.
// Returns nil and an error if the key does not exist in the index.
//
// Example:
//
//	entry, err := idx.Get("sports")
//	if err != nil {
//	    // handle missing key
//	} else {
//	    // use entry.Values
//	}
func (idx *InvertedIndex) Get(key string) (*IndexEntry, error) {
	entry, ok := idx.indexMap[key]
	if !ok {
		zlog.LOG.Warn("InvertedIndex.Get.KeyNotFound", zap.String("key", key))
		return nil, fmt.Errorf("index entry not found for key: %s", key)
	}

	return &entry, nil
}

// GetUpdateTime returns the UNIX timestamp when the index was last updated.
func (idx *InvertedIndex) GetUpdateTime() int64 {
	return idx.updateTime
}

// GetURL returns the source file path of the inverted index.
func (idx *InvertedIndex) GetURL() string {
	return idx.filePath
}
