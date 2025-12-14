package constrains

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// FixedPositionInsert moves an item that matches a given condition to a fixed position in the collection.
// The condition is evaluated via a Minia rule and returns 1 if the entry should be fixed.
type FixedPositionInsert struct {
	conf    *model.FixedPositionInsertedConstrainConfigure // constraint configuration
	program *minia.Minia                                   // compiled condition program
}

// NewFixedPositionInsert creates a fixed-position constraint from configuration.
func NewFixedPositionInsert(conf *model.FixedPositionInsertedConstrainConfigure) *FixedPositionInsert {
	pStat := prome.NewStat("NewFixedPositionInsert")
	defer pStat.End()
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Condition)})
	return &FixedPositionInsert{
		conf:    conf,
		program: program,
	}
}

// Do evaluates all entries using the configured condition.
// If an entry is hit (condition returns 1), it will be moved to f.conf.Position.
// Only the first matching entry is moved; others remain in their positions.
func (f *FixedPositionInsert) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("FixedPositionInsert.Do")
	defer pStat.End()
	if f.conf.Position >= len(collection) {
		// Target position is out of range
		return collection
	}

	for i, entry := range collection {
		// Ensure parameter order matches the rest of the engine: basic, user features, runtime
		value := f.program.Eval(entry.Runtime.Basic, uCtx.Features, entry.Runtime.RunTime)
		result := value.Get("result")
		if result == nil {
			continue
		}

		hit, err := result.GetInt64()
		if err == nil && hit == 1 {
			zlog.LOG.Info("FixedPositionInsert.Hit",
				zap.String("key", entry.KeyScore.Key),
				zap.Int("current_index", i),
				zap.Int("target_position", f.conf.Position))

			if i == f.conf.Position {
				// Already at desired position
				return collection
			}

			// Move entry to target position
			ret := make(model.Collection, 0, len(collection))
			for idx, e := range collection {
				if idx == i {
					continue // skip original position
				}
				if idx == f.conf.Position {
					ret = append(ret, entry) // insert at target
				}
				ret = append(ret, e)
			}
			return ret
		}
	}

	// No matching entry found
	return collection
}
