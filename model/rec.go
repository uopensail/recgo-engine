package model

import (
	"github.com/uopensail/recgo-engine/model/dbmodel/table"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/sample"
)

type IFliter interface {
	Check(id int) bool
}

type RecallRecord struct {
	ID    int
	Score float32
	Info  string
}
type ItemFeatures struct {
	Source     *pool.Features
	MutFeature *sample.MutableFeatures
}

func (f *ItemFeatures) Keys() []string {

	ret := make([]string, 0, f.Source.Feats.Len()+f.MutFeature.Len())
	ret = append(ret, f.Source.Feats.Keys()...)
	ret = append(ret, f.MutFeature.Keys()...)
	return ret
}

func (f *ItemFeatures) Len() int {
	return f.Source.Feats.Len() + f.MutFeature.Len()
}

func (f *ItemFeatures) Get(key string) sample.Feature {
	if v := f.MutFeature.Get(key); v != nil {
		return v
	}
	return f.Source.Feats.Get(key)
}

func (f *ItemFeatures) Set(key string, value sample.Feature) {
	f.MutFeature.Set(key, value)
}

func (f *ItemFeatures) MarshalJSON() ([]byte, error) {
	feats := sample.NewMutableFeatures()
	keys := f.Source.Feats.Keys()
	for _, key := range keys {
		feat := f.Source.Get(key)
		feats.Set(key, feat)
	}

	keys = f.MutFeature.Keys()
	for _, key := range keys {
		feat := f.MutFeature.Get(key)
		feats.Set(key, feat)
	}

	return feats.MarshalJSON()
}

func (f *ItemFeatures) UnmarshalJSON(data []byte) error {
	return f.Source.Feats.UnmarshalJSON(data)
}

type ItemRefScore struct {
	ItemFeatures
	RecallRecord
}
type ItemRefList []*pool.Features

type ItemScoreList []ItemRefScore

func (list ItemScoreList) Less(i, j int) bool {
	return list[i].Score > list[j].Score
}

func (list ItemScoreList) Len() int {
	return len(list)
}

func (list ItemScoreList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type ItemRefListSortHelper struct {
	ItemRefList
	SortFieldName string
}

func (x ItemRefListSortHelper) Len() int { // 重写 Len() 方法
	return len(x.ItemRefList)
}
func (x ItemRefListSortHelper) Swap(i, j int) { // 重写 Swap() 方法
	x.ItemRefList[i], x.ItemRefList[j] = x.ItemRefList[j], x.ItemRefList[i]
}
func (x ItemRefListSortHelper) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	a := x.ItemRefList[i]
	b := x.ItemRefList[j]
	aScore, _ := a.Feats.Get(x.SortFieldName).GetFloat32()
	bScore, _ := b.Feats.Get(x.SortFieldName).GetFloat32()
	return aScore > bScore
}

type RecallResult struct {
	Items ItemScoreList

	Meta *table.RecallEntityMeta
}

type RecItemList ItemScoreList

func (list RecItemList) GetRecCollection() []int {
	ret := make([]int, len(list))
	for i := 0; i < len(list); i++ {
		ret[i] = list[i].ID
	}
	return ret
}

type ItemRecallTrace struct {
	RecallIndex []int //   索引

}

// recall 结束不会再变了
type RecallTrace struct {
	//TODO: use key的编号
	ItemRecallNameMap map[int]*ItemRecallTrace //   索引
	RecallResults     []RecallResult
}

func (rt *RecallTrace) ExistRecall(itemInsideNum int, recallIndex int) bool {
	if v, ok := rt.ItemRecallNameMap[itemInsideNum]; ok {
		for i := 0; i < len(v.RecallIndex); i++ {
			if v.RecallIndex[i] == recallIndex {
				return true
			}
		}
	}
	return false
}
func (rt *RecallTrace) GetRecallNames(itemInsideNum int) []string {
	if v, ok := rt.ItemRecallNameMap[itemInsideNum]; ok {
		recalls := make([]string, len(v.RecallIndex))
		for i := 0; i < len(v.RecallIndex); i++ {
			recalls[i] = rt.RecallResults[v.RecallIndex[i]].Meta.EntityMeta.Name
		}
		return recalls
	}
	return nil
}

func (rt *RecallTrace) GetRecallIndex(name string) int {
	for i := 0; i < len(rt.RecallResults); i++ {
		if rt.RecallResults[i].Meta.EntityMeta.Name == name {
			return i
		}
	}
	return -1
}

type WeighedTrace struct {
	InScoreRecord map[int]float32
	WeighedRecord map[int]string
}

type RankTrace struct {
	ScoreRecord map[int]float32 //打分记录
}

type LayoutTrace struct {
	BreakLayoutRecord map[int]int
}

type StageResult struct {
	StageList RecItemList

	RecallTrace
	WeighedTrace
	RankTrace
	LayoutTrace
}

func MergeStageResult(a StageResult, b StageResult) StageResult {

	c := StageResult{
		StageList:   make(RecItemList, 0, len(a.StageList)+len(b.StageList)),
		RecallTrace: RecallTrace{},
	}
	c.StageList = append(a.StageList, b.StageList...)

	c.ItemRecallNameMap = a.ItemRecallNameMap
	aRecallLen := len(a.RecallResults)
	c.RecallResults = append(a.RecallResults, b.RecallResults...)
	for k, v := range b.ItemRecallNameMap {
		for i := 0; i < len(v.RecallIndex); i++ {
			vi := v.RecallIndex[i] + aRecallLen
			if _, ok := c.ItemRecallNameMap[k]; ok == false {
				c.ItemRecallNameMap[k] = &ItemRecallTrace{}
			}
			c.ItemRecallNameMap[k].RecallIndex = append(c.ItemRecallNameMap[k].RecallIndex, vi)

		}

	}

	return c
}

type StatusResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
