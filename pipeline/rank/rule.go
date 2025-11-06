package rank

import (
	"fmt"
	"sort"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/minia"
)

type Rule struct {
	conf    *model.RuleBasedRankConfigure
	program *minia.Minia
}

func NewRule(conf *model.RuleBasedRankConfigure) *Rule {
	program := minia.NewMinia([]string{fmt.Sprintf("result=%s", conf.Rule)})
	return &Rule{
		conf:    conf,
		program: program,
	}
}

func (rule *Rule) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	for _, entry := range collection {
		result := rule.program.Eval(entry.Runtime.Basic, uCtx.Features, entry.Runtime.RunTime)
		score := result.Get("result")
		entry.KeyScore.Score = 0
		if score != nil {
			val, err := score.GetFloat32()
			if err == nil {
				entry.KeyScore.Score = val
			}
		}
	}

	sort.Stable(collection)
	return collection
}
