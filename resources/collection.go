package resources

import (
	"sort"

	"github.com/uopensail/ulib/datastruct"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/uno"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type Collection []int

// BinarySearch 二分查找函数
func BinarySearch(arr Collection, target int) bool {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := low + (high-low)/2 // 避免溢出
		if arr[mid] == target {
			return true // 找到目标值
		} else if arr[mid] < target {
			low = mid + 1 // 目标值在右半部分
		} else {
			high = mid - 1 // 目标值在左半部分
		}
	}

	return false // 未找到目标值
}

func BuildCollectionBitmap(pl *pool.Pool, collection Collection, sourceName string, condition string) datastruct.BitMap {
	stat := prome.NewStat("BuildCollection")
	defer stat.End()
	evaluator, err := uno.NewEvaluator(condition)
	if err != nil {
		zlog.LOG.Error("build collection condition error", zap.Error(err))
		stat.MarkErr()
		return nil
	}
	defer evaluator.Release()

	var status int32
	ret := datastruct.CreateBitMap(pl.Len())
	for i := 0; i < len(collection); i++ {
		slice := evaluator.Allocate()
		itemID := collection[i]
		item := pl.GetById(itemID)
		if item != nil {
			evaluator.Fill(sourceName, &item.Feats, slice)
			status = evaluator.Eval(slice)
			if status == 1 {
				ret.MarkTrue(itemID)
			}
		}

	}
	stat.SetCounter(len(ret))

	return ret
}

func BuildCollection(pl *pool.Pool, collection Collection, sourceName string, condition string) Collection {
	stat := prome.NewStat("BuildCollection")
	defer stat.End()
	evaluator, err := uno.NewEvaluator(condition)
	if err != nil {
		zlog.LOG.Error("build collection condition error", zap.Error(err))
		stat.MarkErr()
		return nil
	}
	defer evaluator.Release()

	var status int32
	ret := make([]int, 0, len(collection))
	for i := 0; i < len(collection); i++ {
		slice := evaluator.Allocate()
		id := collection[i]
		item := pl.GetById(id)
		if item != nil {
			evaluator.Fill(sourceName, &item.Feats, slice)
			status = evaluator.Eval(slice)
			if status == 1 {
				ret = append(ret, id)
			}
		}

	}
	stat.SetCounter(len(ret))

	return ret
}

func OrderBy(pl *pool.Pool, collection Collection, key string, desc bool) {
	stat := prome.NewStat("Source.Sort")
	defer stat.End()
	if len(key) == 0 {
		stat.MarkErr()
		zlog.LOG.Error("sort key is nil")
		return
	}

	less := func(i, j int) bool {
		left := pl.Array[i].Get(key)
		right := pl.Array[j].Get(key)

		if left == nil || right == nil {
			return false
		}

		dtype := left.Type()
		switch dtype {
		case sample.Float32Type:
			lv, err1 := left.GetFloat32()
			rv, err2 := right.GetFloat32()
			if err1 != nil || err2 != nil {
				return false
			}
			return (desc && lv > rv) || (!desc && lv < rv)

		case sample.Int64Type:
			lv, err1 := left.GetInt64()
			rv, err2 := right.GetInt64()
			if err1 != nil || err2 != nil {
				return false
			}
			return (desc && lv > rv) || (!desc && lv < rv)
		case sample.StringType:
			lv, err1 := left.GetString()
			rv, err2 := right.GetString()
			if err1 != nil || err2 != nil {
				return false
			}
			return (desc && lv > rv) || (!desc && lv < rv)
		default:
			return false
		}
	}
	sort.SliceStable(collection, func(i, j int) bool {
		return less(collection[i], collection[j])
	})
}

func Intersection(max int, collections ...Collection) Collection {

	cl := len(collections)
	if cl == 0 {
		return nil
	}
	if cl == 1 {
		return collections[0]
	}
	a := datastruct.CreateBitMap(max)
	for i := 0; i < len(collections[0]); i++ {
		a.MarkTrue(collections[0][i])
	}
	b := datastruct.CreateBitMap(max)
	for i := 1; i < len(collections); i++ {
		for j := 0; j < len(collections[i]); j++ {
			b.MarkTrue(collections[i][j])
		}
		a.And(b)
		b.Clear()
	}
	ret := make([]int, 0, max/len(collections))
	for i := 0; i < len(a); i++ {
		v := a[i]
		for j := 0; j < 8; j++ {
			if (v & 1) == 1 {
				ret = append(ret, i*8+j)
			}
			v >>= 1
		}
	}
	return ret
}
