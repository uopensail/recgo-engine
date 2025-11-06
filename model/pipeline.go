package model

import (
	"encoding/json"
	"fmt"
)

// ============ Constants ============

const (
	// Recall types
	RecallTypeMatch = "match" ///< Match-based recall type
	RecallTypeModel = "model" ///< Model-based recall type

	// Rank types
	RankTypeRule            = "rule"             ///< Rule-based ranking type
	RankTypeChannelPriority = "channel_priority" ///< Channel priority-based ranking type
	RankTypeModel           = "model"            ///< Model-based ranking type

	// Constraint types
	ConstraintTypeScatter        = "scatter"         ///< Scatter-based constraint type
	ConstraintTypeWeightAdjusted = "weight_adjusted" ///< Weight-adjusted constraint type
	ConstraintTypeFixedPosition  = "fixed_position"  ///< Fixed position constraint type
)

// ============ Interfaces ============

/**
 * @brief Interface for frequency configuration
 */
type IFreq interface {
	GetName() string
	GetTimespan() int
	GetFrequency() int
	GetAction() string
}

/**
 * @brief Interface for recall configuration
 */
type IRecall interface {
	GetName() string
	GetType() string
	GetCount() int
}

/**
 * @brief Interface for ranking configuration
 */
type IRank interface {
	GetName() string
	GetType() string
}

/**
 * @brief Interface for constraint configuration
 */
type IConstrain interface {
	GetName() string
	GetType() string
}

// ============ Frequency Configuration ============

/**
 * @brief Frequency configuration structure
 * @details Defines frequency control parameters for pipeline execution
 */
type FreqConfigure struct {
	Name      string `json:"name"`      ///< Configuration name
	Timespan  int    `json:"timespan"`  ///< Time span in seconds
	Frequency int    `json:"frequency"` ///< Maximum frequency within timespan
	Action    string `json:"action"`    ///< Action to take when limit exceeded
}

/**
 * @brief Get the name of frequency configuration
 * @return Configuration name
 */
func (f FreqConfigure) GetName() string { return f.Name }

/**
 * @brief Get the timespan of frequency configuration
 * @return Timespan in seconds
 */
func (f FreqConfigure) GetTimespan() int { return f.Timespan }

/**
 * @brief Get the frequency limit
 * @return Maximum frequency
 */
func (f FreqConfigure) GetFrequency() int { return f.Frequency }

/**
 * @brief Get the action when frequency limit is exceeded
 * @return Action string
 */
func (f FreqConfigure) GetAction() string { return f.Action }

// ============ Recall Configurations ============

/**
 * @brief Match-based recall configuration
 * @details Uses expression matching for item recall
 */
type MatchRecallConfigure struct {
	Name  string `json:"name"`  ///< Configuration name
	Type  string `json:"type"`  ///< Configuration type (should be RecallTypeMatch)
	Expr  string `json:"expr"`  ///< Match expression
	Index string `json:"index"` ///< Index name to search
	Count int    `json:"count"` ///< Maximum number of items to recall
}

/**
 * @brief Get the name of match recall configuration
 * @return Configuration name
 */
func (m MatchRecallConfigure) GetName() string { return m.Name }

/**
 * @brief Get the type of match recall configuration
 * @return Configuration type
 */
func (m MatchRecallConfigure) GetType() string { return m.Type }

/**
 * @brief Get the count of match recall configuration
 * @return Maximum recall count
 */
func (m MatchRecallConfigure) GetCount() int { return m.Count }

/**
 * @brief Model-based recall configuration
 * @details Uses external model service for item recall
 */
type ModelRecallConfigure struct {
	Name  string `json:"name"`  ///< Configuration name
	Type  string `json:"type"`  ///< Configuration type (should be RecallTypeModel)
	URL   string `json:"url"`   ///< Model service URL
	Count int    `json:"count"` ///< Maximum number of items to recall
}

/**
 * @brief Get the name of model recall configuration
 * @return Configuration name
 */
func (m ModelRecallConfigure) GetName() string { return m.Name }

/**
 * @brief Get the type of model recall configuration
 * @return Configuration type
 */
func (m ModelRecallConfigure) GetType() string { return m.Type }

/**
 * @brief Get the count of model recall configuration
 * @return Maximum recall count
 */
func (m ModelRecallConfigure) GetCount() int { return m.Count }

// ============ Ranking Configurations ============

/**
 * @brief Rule-based ranking configuration
 * @details Uses predefined rules for item ranking
 */
type RuleBasedRankConfigure struct {
	Name string `json:"name"` ///< Configuration name
	Type string `json:"type"` ///< Configuration type (should be RankTypeRule)
	Rule string `json:"rule"` ///< Ranking rule expression
}

/**
 * @brief Get the name of rule-based rank configuration
 * @return Configuration name
 */
func (r RuleBasedRankConfigure) GetName() string { return r.Name }

/**
 * @brief Get the type of rule-based rank configuration
 * @return Configuration type
 */
func (r RuleBasedRankConfigure) GetType() string { return r.Type }

/**
 * @brief Channel priority-based ranking configuration
 * @details Ranks items based on channel priorities
 */
type ChannelPriorityRankConfigure struct {
	Name       string   `json:"name"`       ///< Configuration name
	Type       string   `json:"type"`       ///< Configuration type (should be RankTypeChannelPriority)
	Priorities []string `json:"priorities"` ///< Ordered list of channel priorities
}

/**
 * @brief Get the name of channel priority rank configuration
 * @return Configuration name
 */
func (c ChannelPriorityRankConfigure) GetName() string { return c.Name }

/**
 * @brief Get the type of channel priority rank configuration
 * @return Configuration type
 */
func (c ChannelPriorityRankConfigure) GetType() string { return c.Type }

/**
 * @brief Model-based ranking configuration
 * @details Uses external model service for item ranking
 */
type ModelBasedRankConfigure struct {
	Name string `json:"name"` ///< Configuration name
	Type string `json:"type"` ///< Configuration type (should be RankTypeModel)
	URL  string `json:"url"`  ///< Model service URL
}

/**
 * @brief Get the name of model-based rank configuration
 * @return Configuration name
 */
func (m ModelBasedRankConfigure) GetName() string { return m.Name }

/**
 * @brief Get the type of model-based rank configuration
 * @return Configuration type
 */
func (m ModelBasedRankConfigure) GetType() string { return m.Type }

// ============ Constraint Configurations ============

/**
 * @brief Scatter-based constraint configuration
 * @details Ensures diversity by scattering items based on a field
 */
type ScatterBasedConstrainConfigure struct {
	Name  string `json:"name"`  ///< Configuration name
	Type  string `json:"type"`  ///< Configuration type (should be ConstraintTypeScatter)
	Field string `json:"field"` ///< Field name for scattering
	Count int    `json:"count"` ///< Maximum items per field value
}

/**
 * @brief Get the name of scatter-based constraint configuration
 * @return Configuration name
 */
func (s ScatterBasedConstrainConfigure) GetName() string { return s.Name }

/**
 * @brief Get the type of scatter-based constraint configuration
 * @return Configuration type
 */
func (s ScatterBasedConstrainConfigure) GetType() string { return s.Type }

/**
 * @brief Weight-adjusted constraint configuration
 * @details Adjusts item weights based on conditions
 */
type WeightAdjustedConstrainConfigure struct {
	Name      string  `json:"name"`      ///< Configuration name
	Type      string  `json:"type"`      ///< Configuration type (should be ConstraintTypeWeightAdjusted)
	Ratio     float32 `json:"ratio"`     ///< Weight adjustment ratio
	Condition string  `json:"condition"` ///< Condition for weight adjustment
}

/**
 * @brief Get the name of weight-adjusted constraint configuration
 * @return Configuration name
 */
func (w WeightAdjustedConstrainConfigure) GetName() string { return w.Name }

/**
 * @brief Get the type of weight-adjusted constraint configuration
 * @return Configuration type
 */
func (w WeightAdjustedConstrainConfigure) GetType() string { return w.Type }

/**
 * @brief Fixed position insertion constraint configuration
 * @details Inserts items at fixed positions based on conditions
 */
type FixedPositionInsertedConstrainConfigure struct {
	Name      string `json:"name"`      ///< Configuration name
	Type      string `json:"type"`      ///< Configuration type (should be ConstraintTypeFixedPosition)
	Position  int    `json:"position"`  ///< Fixed position for insertion
	Condition string `json:"condition"` ///< Condition for insertion
}

/**
 * @brief Get the name of fixed position constraint configuration
 * @return Configuration name
 */
func (f FixedPositionInsertedConstrainConfigure) GetName() string { return f.Name }

/**
 * @brief Get the type of fixed position constraint configuration
 * @return Configuration type
 */
func (f FixedPositionInsertedConstrainConfigure) GetType() string { return f.Type }

// ============ Pipeline Configuration ============

/**
 * @brief Complete pipeline configuration
 * @details Defines the entire recommendation pipeline including frequency control,
 *          recall strategies, ranking methods, and constraints
 */
type PipelineConfigure struct {
	Name       string       `json:"name"`           ///< Pipeline name
	Freqs      []IFreq      `json:"freqs"`          ///< Frequency configurations
	Recalls    []IRecall    `json:"recalls"`        ///< Recall configurations
	Rank       IRank        `json:"rank,omitempty"` ///< Ranking configuration (optional)
	Constrains []IConstrain `json:"constrains"`     ///< Constraint configurations
	Buckets    []int        `json:"bucktes"`        ///< Hit buckets
}

/**
 * @brief Custom JSON unmarshaling for PipelineConfigure
 * @param data JSON data to unmarshal
 * @return Error if unmarshaling fails
 */
func (p *PipelineConfigure) UnmarshalJSON(data []byte) error {
	// Define temporary structure for parsing
	var temp struct {
		Name       string            `json:"name"`
		Freqs      []FreqConfigure   `json:"freqs"`
		Recalls    []json.RawMessage `json:"recalls"`
		Rank       json.RawMessage   `json:"rank,omitempty"`
		Constrains []json.RawMessage `json:"constrains"`
		Buckets    []int             `json:"bucktes"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal pipeline configure: %w", err)
	}

	// Set basic fields
	p.Name = temp.Name
	p.Buckets = temp.Buckets
	p.Freqs = make([]IFreq, 0, len(temp.Freqs))
	for _, freq := range temp.Freqs {
		p.Freqs = append(p.Freqs, freq)
	}

	// Process recalls
	if err := p.unmarshalRecalls(temp.Recalls); err != nil {
		return err
	}

	// Process rank (optional)
	if len(temp.Rank) > 0 {
		if err := p.unmarshalRank(temp.Rank); err != nil {
			return err
		}
	}

	// Process constraints
	if err := p.unmarshalConstrains(temp.Constrains); err != nil {
		return err
	}

	return nil
}

/**
 * @brief Unmarshal recalls from JSON raw messages
 * @param rawRecalls Array of JSON raw messages
 * @return Error if unmarshaling fails
 */
func (p *PipelineConfigure) unmarshalRecalls(rawRecalls []json.RawMessage) error {
	p.Recalls = make([]IRecall, len(rawRecalls))

	for i, raw := range rawRecalls {
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeCheck); err != nil {
			return fmt.Errorf("failed to get recall type at index %d: %w", i, err)
		}

		switch typeCheck.Type {
		case RecallTypeMatch:
			var config MatchRecallConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal match recall at index %d: %w", i, err)
			}
			p.Recalls[i] = config
		case RecallTypeModel:
			var config ModelRecallConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal model recall at index %d: %w", i, err)
			}
			p.Recalls[i] = config
		default:
			return fmt.Errorf("unknown recall type '%s' at index %d", typeCheck.Type, i)
		}
	}

	return nil
}

/**
 * @brief Unmarshal rank from JSON raw message
 * @param rawRank JSON raw message
 * @return Error if unmarshaling fails
 */
func (p *PipelineConfigure) unmarshalRank(rawRank json.RawMessage) error {
	var typeCheck struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(rawRank, &typeCheck); err != nil {
		return fmt.Errorf("failed to get rank type: %w", err)
	}

	switch typeCheck.Type {
	case RankTypeRule:
		var config RuleBasedRankConfigure
		if err := json.Unmarshal(rawRank, &config); err != nil {
			return fmt.Errorf("failed to unmarshal rule rank: %w", err)
		}
		p.Rank = config
	case RankTypeChannelPriority:
		var config ChannelPriorityRankConfigure
		if err := json.Unmarshal(rawRank, &config); err != nil {
			return fmt.Errorf("failed to unmarshal channel priority rank: %w", err)
		}
		p.Rank = config
	case RankTypeModel:
		var config ModelBasedRankConfigure
		if err := json.Unmarshal(rawRank, &config); err != nil {
			return fmt.Errorf("failed to unmarshal model rank: %w", err)
		}
		p.Rank = config
	default:
		return fmt.Errorf("unknown rank type: %s", typeCheck.Type)
	}

	return nil
}

/**
 * @brief Unmarshal constraints from JSON raw messages
 * @param rawConstrains Array of JSON raw messages
 * @return Error if unmarshaling fails
 */
func (p *PipelineConfigure) unmarshalConstrains(rawConstrains []json.RawMessage) error {
	p.Constrains = make([]IConstrain, len(rawConstrains))

	for i, raw := range rawConstrains {
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeCheck); err != nil {
			return fmt.Errorf("failed to get constraint type at index %d: %w", i, err)
		}

		switch typeCheck.Type {
		case ConstraintTypeScatter:
			var config ScatterBasedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal scatter constraint at index %d: %w", i, err)
			}
			p.Constrains[i] = config
		case ConstraintTypeWeightAdjusted:
			var config WeightAdjustedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal weight adjusted constraint at index %d: %w", i, err)
			}
			p.Constrains[i] = config
		case ConstraintTypeFixedPosition:
			var config FixedPositionInsertedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal fixed position constraint at index %d: %w", i, err)
			}
			p.Constrains[i] = config
		default:
			return fmt.Errorf("unknown constraint type '%s' at index %d", typeCheck.Type, i)
		}
	}

	return nil
}
