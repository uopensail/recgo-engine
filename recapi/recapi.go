package recapi

import (
	"github.com/uopensail/ulib/sample"
)

// Request represents the recommendation API request parameters.
type Request struct {
	TraceId  string                  `json:"trace_id,omitempty"`  // Unique request ID for tracking
	UserId   string                  `json:"user_id"`             // User ID
	Pipeline string                  `json:"pipeline"`            // Pipeline name or ID to execute
	RelateId string                  `json:"relate_id,omitempty"` // Related content ID (e.g., product ID in a detail page)
	Count    int32                   `json:"count,omitempty"`     // Number of items to request
	Context  *sample.MutableFeatures `json:"context,omitempty"`   // Contextual features (e.g., session info)
	Features *sample.MutableFeatures `json:"features,omitempty"`  // External features provided by caller
}

// ItemInfo represents a single recommended item with channels and reasons.
type ItemInfo struct {
	Item     string   `json:"item,omitempty"`     // Item ID
	Channels []string `json:"channels,omitempty"` // Channels that contributed to recommendation
	Reasons  []string `json:"reasons,omitempty"`  // Reasons or explanations for recommendation
}

// Response represents the standard API response for recommendation results.
type Response struct {
	Code     int         `json:"code"`               // Business status code: 0 = success, non-zero = error
	Message  string      `json:"message,omitempty"`  // Message for status or error description
	TraceId  string      `json:"trace_id,omitempty"` // Request ID
	UserId   string      `json:"user_id,omitempty"`  // User ID
	Pipeline string      `json:"pipeline,omitempty"` // Pipeline name used
	Items    []*ItemInfo `json:"items,omitempty"`    // Recommended items
	Count    int         `json:"count,omitempty"`    // Number of items returned
}
