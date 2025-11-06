package model

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

/**
 * @struct KeyScore
 * @brief Represents a id item with its relevance score
 */
type KeyScore struct {
	Key   string  `json:"key"`   ///< Unique identifier of the candidate item
	Score float32 `json:"score"` ///< Relevance score for the candidate item
}

/**
 * @struct IndexEntry
 * @brief Represents a single entry in the inverted index
 */
type IndexEntry struct {
	Key    string     `json:"key"`    ///< Index key (feature value)
	Values []KeyScore `json:"values"` ///< List of candidates associated with this key
}

/**
 * @struct InvertedIndex
 * @brief Inverted index for fast candidate retrieval
 */
type InvertedIndex struct {
	Indexes    map[string]IndexEntry ///< Internal map for O(1) key lookup
	filePath   string                ///< Filepath of inverted index
	updatetime int64
}

/**
 * @brief Creates a new InvertedIndex instance from a JSON file
 * @param filePath Path to the JSON file containing index data
 * @return Pointer to a new InvertedIndex, Error if loading fails
 *
 * The JSON file should contain an array of IndexEntry objects.
 * Candidates within each entry will be sorted by score in descending order.
 */
func NewInvertedIndex(filePath string) (Resource, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse JSON data
	var entries []IndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from file %s: %w", filePath, err)
	}

	// Build index map
	indexes := make(map[string]IndexEntry, len(entries))
	for _, entry := range entries {
		// Sort candidates by score in descending order
		sort.Slice(entry.Values, func(i, j int) bool {
			return entry.Values[i].Score > entry.Values[j].Score
		})
		indexes[entry.Key] = entry
	}

	return &InvertedIndex{Indexes: indexes, filePath: filePath, updatetime: time.Now().Unix()}, nil
}

func (idx *InvertedIndex) GetUpdateTime() (int64, error) {
	return idx.updatetime, nil
}

func (idx *InvertedIndex) GetURL() (string, error) {
	return idx.filePath, nil
}
