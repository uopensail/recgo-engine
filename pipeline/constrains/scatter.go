package constrains

import (
	"strconv"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Scatter enforces distribution constraints on a collection of entries.
// Each ScatterBasedConstrainConfigure specifies a feature field and a maximum allowed count.
// Entries exceeding the allowed count for a feature value are deferred to the end of the collection.
type Scatter struct {
	confs []*model.ScatterBasedConstrainConfigure
}

// NewScatter creates a new Scatter constraint handler using the provided configurations.
func NewScatter(confs []*model.ScatterBasedConstrainConfigure) *Scatter {
	pStat := prome.NewStat("NewScatter")
	defer pStat.End()
	return &Scatter{
		confs: confs,
	}
}

// Do applies scatter constraints:
// 1. For each configured field, maintain a count of how many times each value appears in the accepted list.
// 2. An entry is accepted only if none of its configured feature values exceed the count limit.
// 3. Entries violating any constraint are moved to the end of the collection.
func (s *Scatter) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	pStat := prome.NewStat("Scatter.Do")
	defer pStat.End()
	// filter[i] holds value counts for conf[i]
	filter := make([]map[string]int, len(s.confs))
	for i := range s.confs {
		filter[i] = make(map[string]int)
	}

	ret := make([]*model.Entry, 0, len(collection))    // accepted entries
	remain := make([]*model.Entry, 0, len(collection)) // violating entries

	for _, entry := range collection {
		status := true // entry passes all constraints?

		// Local keys slice for current entry
		keys := make([][]string, len(s.confs))

		// Check each scatter config
		for i := 0; i < len(s.confs) && status; i++ {
			conf := s.confs[i]
			fea, err := entry.Get(conf.Field)
			if err != nil {
				// Feature not found in this entry, skip check
				continue
			}

			// Convert feature value(s) to string slice
			keys[i] = Feature2StringSlice(fea)

			// Check violation for any of the feature keys
			for _, key := range keys[i] {
				if filter[i][key] >= conf.Count {
					status = false
					zlog.LOG.Debug("Scatter.ConstraintViolated",
						zap.String("entry_key", entry.KeyScore.Key),
						zap.String("field", conf.Field),
						zap.String("violating_value", key),
						zap.Int("allowed_count", conf.Count),
						zap.Int("current_count", filter[i][key]))
					break
				}
			}
		}

		if status {
			// Accept entry and update counts
			ret = append(ret, entry)
			for i := 0; i < len(s.confs); i++ {
				for _, key := range keys[i] {
					filter[i][key]++
					zlog.LOG.Debug("Scatter.UpdateCount",
						zap.String("field", s.confs[i].Field),
						zap.String("value", key),
						zap.Int("new_count", filter[i][key]))
				}
			}
		} else {
			// Violating entry goes to remain list
			remain = append(remain, entry)
		}
	}

	// Append violating entries at the end
	ret = append(ret, remain...)

	zlog.LOG.Debug("Scatter.Completed",
		zap.Int("total_entries", len(collection)),
		zap.Int("accepted_entries", len(ret)-len(remain)),
		zap.Int("violated_entries", len(remain)))

	return ret
}

// Feature2StringSlice converts a sample.Feature value into a slice of strings.
// This is required because scatter constraints are applied on string keys.
func Feature2StringSlice(feature sample.Feature) []string {
	switch feature.Type() {
	case sample.Int64Type:
		val, _ := feature.GetInt64()
		return []string{strconv.FormatInt(val, 10)}
	case sample.Int64sType:
		val, _ := feature.GetInt64s()
		ret := make([]string, 0, len(val))
		for _, v := range val {
			ret = append(ret, strconv.FormatInt(v, 10))
		}
		return ret
	case sample.StringType:
		val, _ := feature.GetString()
		return []string{val}
	case sample.StringsType:
		val, _ := feature.GetStrings()
		return val
	default:
		return nil
	}
}
