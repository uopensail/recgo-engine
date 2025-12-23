package recalls

import (
	"fmt"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/program"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

// Matcher implements match-based recall using Minia expressions and an InvertedIndex.
type Matcher struct {
	conf    *model.MatchRecallConfigure // Recall configuration
	index   *model.InvertedIndex        // Candidate index
	program *program.Program            // Compiled expression program
}

// NewMatcher creates a new Matcher instance with compiled Minia expression.
func NewMatcher(conf *model.MatchRecallConfigure) *Matcher {
	pStat := prome.NewStat("Recall.NewMatcher")
	defer pStat.End()

	program, err := program.NewProgram(conf.Expr)
	if err != nil {
		zlog.LOG.Error("Recall.NewMatcher program create error",
			zap.Error(err))
		panic(err)
	}

	zlog.LOG.Info("Matcher.Created", zap.String("name", conf.Name), zap.String("expr", conf.Expr))
	return &Matcher{
		conf:    conf,
		index:   nil, // needs to be injected externally
		program: program,
	}
}

// Do runs the recall process and returns candidate entries.
// 1. Evaluate expression to get matching keys.
// 2. Merge candidates across keys using layered round-robin.
// 3. Return top N entries based on conf.Count.
func (m *Matcher) Do(uCtx *userctx.UserContext) model.Collection {
	pStat := prome.NewStat("Recall.Matcher.Do")
	defer pStat.End()

	if m.index == nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Matcher.Do.NoIndex", zap.String("recall", m.conf.Name))
		return nil
	}

	var value any
	var err error
	if uCtx.Related != nil {
		value, err = m.program.Eval(uCtx.Related, uCtx.Features)
	} else {
		value, err = m.program.Eval(uCtx.Features)
	}

	if err != nil {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Matcher.Do program eval value error", zap.Error(err))
		return nil
	}

	keys, ok := value.([]string)
	if !ok {
		pStat.MarkErr()
		zlog.LOG.Error("Recall.Matcher.Do program eval value to []string error", zap.String("recall", m.conf.Name))
		return nil
	}

	zlog.LOG.Debug("Matcher.Do.KeysExtracted", zap.Int("count", len(keys)))

	ret := mergeCandidatesRoundRobin(keys, m.index, uCtx.Items, m.conf.Name)

	if len(ret) == 0 {
		pStat.MarkErr()
		zlog.LOG.Info("Recall.Matcher.Do.NoResult", zap.String("recall", m.conf.Name))
		return nil
	}

	count := minInt(m.conf.Count, len(ret))
	zlog.LOG.Debug("Matcher.Do.Completed",
		zap.Int("requested_count", m.conf.Count),
		zap.Int("returned_count", count),
	)
	pStat.SetCounter(count)
	return ret[:count]
}

// mergeCandidatesRoundRobin merges multiple candidate lists using layered round-robin (Z-like interleaving).
func mergeCandidatesRoundRobin(
	keys []string,
	index *model.InvertedIndex,
	items *model.Items,
	recallName string,
) []*model.Entry {
	filter := make(map[int]struct{})
	ret := make([]*model.Entry, 0)

	// 计算候选最大长度
	maxSize := 0
	for _, key := range keys {
		if entry, err := index.Get(key); err == nil {
			maxSize = maxInt(maxSize, len(entry.Values))
		}
	}

	// 分层轮询合并
	for i := 0; i < maxSize; i++ {
		for _, key := range keys {
			if arr, err := index.Get(key); err == nil {
				if i < len(arr.Values) {
					k := arr.Values[i]
					if id, _ := items.GetByKey(k.Key); id >= 0 {
						if _, exists := filter[id]; exists {
							continue
						}
						filter[id] = struct{}{}
						entry, err := model.NewEntry(k, items)
						if err != nil {
							zlog.LOG.Warn("mergeCandidatesRoundRobin.NewEntryError",
								zap.String("key", k.Key), zap.Error(err))
							continue
						}
						entry.AddChan(recallName, fmt.Sprintf("recall by key: %s", key))
						ret = append(ret, entry)
					}
				}
			}
		}
	}

	return ret
}

// maxInt returns the larger of a and b.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt returns the smaller of a and b.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
