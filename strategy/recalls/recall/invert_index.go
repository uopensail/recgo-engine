package recall

import (
	"math"
	"strconv"
	"strings"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/model/dbmodel"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/recgo-engine/resources"

	"github.com/uopensail/recgo-engine/userctx"

	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/utils"
)

type InvertInexRecall struct {
	table.RecallEntityMeta
	table.InvertInexRecallMeta
	condition *resources.Condition
	*pool.Pool
}

func NewInvertInexRecall(meta table.RecallEntityMeta, ress *resources.Resource, dbModel *dbmodel.DBTabelModel) IRecallStrategyEntity {
	//

	recall := InvertInexRecall{
		Pool:             ress.Pool,
		RecallEntityMeta: meta,
	}
	recall.InvertInexRecallMeta = meta.ParseInvertInexRecallMeta()
	recall.condition = resources.BuildCondition(ress, ress.Pool.WholeCollection, recall.Condition)

	return &recall
}
func (r *InvertInexRecall) Meta() *table.RecallEntityMeta {
	return &r.RecallEntityMeta
}
func (r *InvertInexRecall) Do(uCtx *userctx.UserContext, ifilter model.IFliter) ([]int, error) {
	stat := prome.NewStat("Recall." + r.EntityMeta.Name)
	defer stat.End()
	ret := r.do(uCtx, ifilter)
	stat.SetCounter(len(ret))
	return ret, nil
}
func (r *InvertInexRecall) formatKeys(userFeats sample.Features) []string {

	// Step 1: 按字段顺序收集值列表
	allValues := make([][]string, 0, len(r.UserFeatureFields))
	for _, key := range r.UserFeatureFields {
		fieldValues := userFeats.Get(key)
		var ss []string
		switch fieldValues.Type() {
		case sample.Int64Type:
			v, _ := fieldValues.GetInt64()
			ss = append(ss, utils.Int642String(v))
		case sample.Int64sType:
			vv, _ := fieldValues.GetInt64s()
			for _, v := range vv {
				ss = append(ss, utils.Int642String(v))
			}
		case sample.StringType:
			v, _ := fieldValues.GetString()
			ss = append(ss, v)
		case sample.StringsType:
			vv, _ := fieldValues.GetStrings()
			ss = vv
		case sample.Float32Type, sample.Float32sType:
			continue
		}
		allValues = append(allValues, ss)
	}

	// Step 2: 生成笛卡尔积（支持任意阶数）
	result := make([]string, 0)
	for i, values := range allValues {
		if len(result) == 0 {
			// 初始化第一个字段的值
			for _, v := range values {

				result = append(result, strconv.Itoa(i)+":"+v)
			}
		} else {
			// 迭代生成后续组合
			var temp []string
			for _, existing := range result {
				for _, newVal := range values {
					var builder strings.Builder
					builder.WriteString(existing) // 写入已有部分
					builder.WriteString("|")      // 添加分隔符

					builder.WriteString(strconv.Itoa(i))
					builder.WriteString(":")
					builder.WriteString(newVal) // 追加新值
					temp = append(temp, builder.String())
				}
			}
			result = temp
		}
	}
	return result
}

func zigzagMerge(tmpCollectionList []resources.Collection, eachMaxCol int, ifilter model.IFliter) resources.Collection {
	// 预计算总元素数量
	if len(tmpCollectionList) == 0 {
		return []int{}
	}

	// 预计算总元素数量
	total := 0
	maxCols := 0
	for _, row := range tmpCollectionList {
		total += len(row)
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// 预分配内存
	result := make([]int, 0, total)
	if eachMaxCol <= 0 {
		eachMaxCol = math.MaxInt
	}
	rowCnt := make([]int, len(tmpCollectionList))
	for i := 0; i < maxCols; i++ {
		for j, row := range tmpCollectionList {
			if i < len(row) && ifilter.Check(row[i]) && rowCnt[j] < eachMaxCol {
				result = append(result, row[i])
				rowCnt[j]++
			}
		}
	}

	return result
}
func (r *InvertInexRecall) do(uCtx *userctx.UserContext, filter model.IFliter) []int {

	keys := r.formatKeys(uCtx.UFeat)
	if len(keys) == 0 {
		return nil
	}
	invertIndexCollection, ok := uCtx.Ress.InvertIndexRess[r.Resource]
	if !ok {
		return nil
	}
	tmpCollectionList := invertIndexCollection.Get(keys)
	if len(tmpCollectionList) <= 0 {
		return nil
	}
	//zigzip merger
	tmpCollection := zigzagMerge(tmpCollectionList, r.EachMaxCol, filter)
	if r.TopK > 0 {
		topk := r.TopK
		if len(tmpCollection) > topk {
			return tmpCollection[:topk]
		}
	}

	if r.condition != nil {
		tmpCollection = r.condition.Check(uCtx.UFeat, tmpCollection)
	}
	return tmpCollection
}

func (r *InvertInexRecall) Close() {

}
