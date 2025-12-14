package model

import (
	"encoding/json"
	"fmt"

	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

//
// ================= Constants =================
//

// Recall type constants define methods for retrieving candidate items.
const (
	RecallTypeMatch = "match" // Match-based recall
	RecallTypeModel = "model" // Model-based recall
)

// Rank type constants define methods for ordering candidate items.
const (
	RankTypeRule            = "rule"             // Rule-based ranking
	RankTypeChannelPriority = "channel_priority" // Ranking based on channel priorities
	RankTypeModel           = "model"            // Model-based ranking
)

// Constraint type constants define rules applied to the result set.
const (
	ConstraintTypeScatter        = "scatter"         // Scatter-based constraint
	ConstraintTypeWeightAdjusted = "weight_adjusted" // Adjust weights based on conditions
	ConstraintTypeFixedPosition  = "fixed_position"  // Insert at fixed positions
)

//
// ================= Interfaces =================
//

// IFreq defines the frequency control configuration interface.
type IFreq interface {
	GetName() string
	GetTimespan() int
	GetFrequency() int
	GetAction() string
}

// IRecall defines the recall strategy configuration interface.
type IRecall interface {
	GetName() string
	GetType() string
	GetCount() int
}

// IRank defines the ranking strategy configuration interface.
type IRank interface {
	GetName() string
	GetType() string
}

// IConstrain defines the constraint configuration interface.
type IConstrain interface {
	GetName() string
	GetType() string
}

//
// ================= Frequency Configuration =================
//

// FreqConfigure defines frequency control parameters for pipeline execution.
type FreqConfigure struct {
	Name      string `json:"name"`
	Timespan  int    `json:"timespan"`
	Frequency int    `json:"frequency"`
	Action    string `json:"action"`
}

func (f FreqConfigure) GetName() string   { return f.Name }
func (f FreqConfigure) GetTimespan() int  { return f.Timespan }
func (f FreqConfigure) GetFrequency() int { return f.Frequency }
func (f FreqConfigure) GetAction() string { return f.Action }

//
// ================= Recall Configurations =================
//

// MatchRecallConfigure uses expression matching for item recall.
type MatchRecallConfigure struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Expr  string `json:"expr"`
	Index string `json:"index"`
	Count int    `json:"count"`
}

func (m MatchRecallConfigure) GetName() string { return m.Name }
func (m MatchRecallConfigure) GetType() string { return m.Type }
func (m MatchRecallConfigure) GetCount() int   { return m.Count }

// ModelRecallConfigure uses an external model service for item recall.
type ModelRecallConfigure struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

func (m ModelRecallConfigure) GetName() string { return m.Name }
func (m ModelRecallConfigure) GetType() string { return m.Type }
func (m ModelRecallConfigure) GetCount() int   { return m.Count }

//
// ================= Ranking Configurations =================
//

// RuleBasedRankConfigure uses predefined rules for item ranking.
type RuleBasedRankConfigure struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Rule string `json:"rule"`
}

func (r RuleBasedRankConfigure) GetName() string { return r.Name }
func (r RuleBasedRankConfigure) GetType() string { return r.Type }

// ChannelPriorityRankConfigure ranks items based on channel priorities.
type ChannelPriorityRankConfigure struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Priorities []string `json:"priorities"`
}

func (c ChannelPriorityRankConfigure) GetName() string { return c.Name }
func (c ChannelPriorityRankConfigure) GetType() string { return c.Type }

// ModelBasedRankConfigure uses an external model for ranking items.
type ModelBasedRankConfigure struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (m ModelBasedRankConfigure) GetName() string { return m.Name }
func (m ModelBasedRankConfigure) GetType() string { return m.Type }

//
// ================= Constraint Configurations =================
//

// ScatterBasedConstrainConfigure ensures diversity by scattering items based on a field.
type ScatterBasedConstrainConfigure struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Field string `json:"field"`
	Count int    `json:"count"`
}

func (s ScatterBasedConstrainConfigure) GetName() string { return s.Name }
func (s ScatterBasedConstrainConfigure) GetType() string { return s.Type }

// WeightAdjustedConstrainConfigure adjusts item weights based on conditions.
type WeightAdjustedConstrainConfigure struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Ratio     float32 `json:"ratio"`
	Condition string  `json:"condition"`
}

func (w WeightAdjustedConstrainConfigure) GetName() string { return w.Name }
func (w WeightAdjustedConstrainConfigure) GetType() string { return w.Type }

// FixedPositionInsertedConstrainConfigure inserts items at fixed positions when conditions are met.
type FixedPositionInsertedConstrainConfigure struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Position  int    `json:"position"`
	Condition string `json:"condition"`
}

func (f FixedPositionInsertedConstrainConfigure) GetName() string { return f.Name }
func (f FixedPositionInsertedConstrainConfigure) GetType() string { return f.Type }

//
// ================= Pipeline Configuration =================
//

// PipelineConfigure holds the entire recommendation pipeline configuration.
type PipelineConfigure struct {
	Name       string       `json:"name"`
	Freqs      []IFreq      `json:"freqs"`
	Recalls    []IRecall    `json:"recalls"`
	Rank       IRank        `json:"rank,omitempty"`
	Constrains []IConstrain `json:"constrains"`
}

// UnmarshalJSON customizes JSON decoding for PipelineConfigure.
func (p *PipelineConfigure) UnmarshalJSON(data []byte) error {
	var temp struct {
		Name       string            `json:"name"`
		Freqs      []FreqConfigure   `json:"freqs"`
		Recalls    []json.RawMessage `json:"recalls"`
		Rank       json.RawMessage   `json:"rank,omitempty"`
		Constrains []json.RawMessage `json:"constrains"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		zlog.LOG.Error("PipelineConfigure.UnmarshalError", zap.Error(err))
		return fmt.Errorf("failed to unmarshal pipeline configure: %w", err)
	}

	p.Name = temp.Name

	// Frequency configs
	p.Freqs = make([]IFreq, 0, len(temp.Freqs))
	for _, freq := range temp.Freqs {
		p.Freqs = append(p.Freqs, freq)
	}

	// Recalls
	if err := p.unmarshalRecalls(temp.Recalls); err != nil {
		return err
	}
	zlog.LOG.Info("PipelineConfigure.RecallsLoaded", zap.Int("count", len(p.Recalls)))

	// Rank
	if len(temp.Rank) > 0 {
		if err := p.unmarshalRank(temp.Rank); err != nil {
			return err
		}
		zlog.LOG.Info("PipelineConfigure.RankLoaded", zap.String("type", p.Rank.GetType()))
	}

	// Constraints
	if err := p.unmarshalConstrains(temp.Constrains); err != nil {
		return err
	}
	zlog.LOG.Info("PipelineConfigure.ConstraintsLoaded", zap.Int("count", len(p.Constrains)))

	return nil
}

func (p *PipelineConfigure) unmarshalRecalls(rawRecalls []json.RawMessage) error {
	p.Recalls = make([]IRecall, len(rawRecalls))
	for i, raw := range rawRecalls {
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeCheck); err != nil {
			zlog.LOG.Error("PipelineConfigure.UnmarshalRecalls.TypeError", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to get recall type at index %d: %w", i, err)
		}

		switch typeCheck.Type {
		case RecallTypeMatch:
			var config MatchRecallConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				zlog.LOG.Error("PipelineConfigure.UnmarshalRecalls.MatchError", zap.Int("index", i), zap.Error(err))
				return fmt.Errorf("failed to unmarshal match recall at index %d: %w", i, err)
			}
			p.Recalls[i] = config
		case RecallTypeModel:
			var config ModelRecallConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				zlog.LOG.Error("PipelineConfigure.UnmarshalRecalls.ModelError", zap.Int("index", i), zap.Error(err))
				return fmt.Errorf("failed to unmarshal model recall at index %d: %w", i, err)
			}
			p.Recalls[i] = config
		default:
			return fmt.Errorf("unknown recall type '%s' at index %d", typeCheck.Type, i)
		}
	}
	return nil
}

func (p *PipelineConfigure) unmarshalRank(rawRank json.RawMessage) error {
	var typeCheck struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(rawRank, &typeCheck); err != nil {
		zlog.LOG.Error("PipelineConfigure.UnmarshalRank.TypeError", zap.Error(err))
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

func (p *PipelineConfigure) unmarshalConstrains(rawConstrains []json.RawMessage) error {
	p.Constrains = make([]IConstrain, len(rawConstrains))
	for i, raw := range rawConstrains {
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeCheck); err != nil {
			zlog.LOG.Error("PipelineConfigure.UnmarshalConstraints.TypeError", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to get constraint type at index %d: %w", i, err)
		}

		switch typeCheck.Type {
		case ConstraintTypeScatter:
			var config ScatterBasedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal scatter constraint: %w", err)
			}
			p.Constrains[i] = config
		case ConstraintTypeWeightAdjusted:
			var config WeightAdjustedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal weight adjusted constraint: %w", err)
			}
			p.Constrains[i] = config
		case ConstraintTypeFixedPosition:
			var config FixedPositionInsertedConstrainConfigure
			if err := json.Unmarshal(raw, &config); err != nil {
				return fmt.Errorf("failed to unmarshal fixed position constraint: %w", err)
			}
			p.Constrains[i] = config
		default:
			return fmt.Errorf("unknown constraint type '%s' at index %d", typeCheck.Type, i)
		}
	}
	return nil
}
